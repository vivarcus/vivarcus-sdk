#!/usr/bin/env python3
"""Verify examples demos (dispatched by DEMO_SCENARIO)."""
import os
import sys

from verify_utils import DemoClient, die  # noqa: E402


def req(name: str) -> str:
    val = os.environ.get(name)
    if not val:
        die(f"{name} required")
    return val


def verify_user_action() -> None:
    obj = req("DEMO_OBJECT")
    action = req("DEMO_OBJECT_ACTION")
    field = os.environ.get("DEMO_TITLE_FIELD", "title__c")
    want = req("DEMO_EXPECTED_TITLE")
    client = DemoClient()

    print("=== create record ===")
    body = client.api(
        "POST",
        f"/api/v1/objects/{obj}/records",
        {"fields": {"name__v": "sdk-demo", field: "before"}},
    )
    rec = body.get("record_id") or body.get("id")
    if not rec:
        die(f"create missing record_id: {body}")

    print(f"=== execute Objectaction {action} ===")
    client.api(
        "POST",
        f"/api/v1/objects/{obj}/records/{rec}/actions/execute",
        {"action": action},
    )

    data = client.vivarcus("object", "get", obj, rec).get("data") or {}
    got = data.get(field)
    print(f"  {field}={got!r}")
    if got != want:
        die(f"want {want!r} got {got!r}")
    print("PASS user action updated field")
    note = os.environ.get("DEMO_VERIFY_NOTE", "").strip()
    if note:
        print(note)


def verify_lifecycle() -> None:
    obj = req("DEMO_OBJECT")
    wf_action = req("DEMO_WF_ACTION")
    wf_cancel = req("DEMO_WF_CANCEL")
    field = os.environ.get("DEMO_TITLE_FIELD", "title__c")
    want = req("DEMO_EXPECTED_TITLE")
    submit = os.environ.get("DEMO_SUBMIT_ACTION", "submit__c")
    client = DemoClient()

    def create_record(title: str = "before") -> str:
        body = client.api(
            "POST",
            f"/api/v1/objects/{obj}/records",
            {"fields": {"name__v": "lc-demo", field: title}},
        )
        rec = body.get("record_id") or body.get("id")
        if not rec:
            die(f"create missing record_id: {body}")
        return rec

    def get_title(rec_id: str) -> tuple[object, object]:
        data = client.vivarcus("object", "get", obj, rec_id).get("data") or {}
        return data.get(field), data.get("state__v")

    def assert_title(rec_id: str, label: str) -> None:
        title, state = get_title(rec_id)
        print(f"  {label}: {field}={title!r} state={state!r}")
        if title != want:
            die(f"{label}: want {want!r} got {title!r} state={state!r}")
        print(f"PASS {label}")

    def start_workflow(rec_id: str, wf_name: str) -> tuple[object, object]:
        body = client.api(
            "POST",
            f"/api/v1/objects/{obj}/records/{rec_id}/lifecycle/workflows/start",
            {"workflow": wf_name},
        )
        return body.get("workflow_instance_id"), body.get("workflow_task_id")

    def cancel_workflow(instance_id: str) -> None:
        client.api(
            "POST",
            f"/api/v1/workflow-instances/{instance_id}/cancel",
            {"comment": "demo verify"},
        )

    print("=== verify event_action (create_record) ===")
    rec_event = create_record("before-event")
    assert_title(rec_event, "event_action")

    print("=== verify entry_action (transition → in_review) ===")
    rec_entry = create_record("before-entry")
    title, _ = get_title(rec_entry)
    if title != want:
        die(f"entry setup: create_record event should stamp {field}, got {title!r}")
    client.vivarcus("lifecycle", "transition", obj, rec_entry, "--action", submit)
    assert_title(rec_entry, "entry_action")

    print("=== verify workflow action step ===")
    rec_wf = create_record("before-wf")
    inst, _ = start_workflow(rec_wf, wf_action)
    print(f"  workflow instance {inst}")
    assert_title(rec_wf, "workflow action step")

    print("=== verify workflow cancel ===")
    rec_cancel = create_record("before-cancel")
    inst, task = start_workflow(rec_cancel, wf_cancel)
    print(f"  workflow instance {inst} task {task}")
    if not inst:
        die("cancel test: missing workflow_instance_id")
    if not task:
        die("cancel test: expected active usertask before cancel")
    cancel_workflow(inst)
    assert_title(rec_cancel, "workflow cancel")

    print("\nALL PASS: event / entry / workflow step / cancel")


def verify_multi_component() -> None:
    obj = req("DEMO_OBJECT")
    field = os.environ.get("DEMO_TITLE_FIELD", "title__c")
    set_action = req("DEMO_SET_TITLE_ACTION")
    clear_action = req("DEMO_CLEAR_TITLE_ACTION")
    want_set = req("DEMO_SET_TITLE_EXPECTED")
    want_clear = os.environ.get("DEMO_CLEAR_TITLE_EXPECTED", "")
    client = DemoClient()

    print("=== create record ===")
    body = client.api(
        "POST",
        f"/api/v1/objects/{obj}/records",
        {"fields": {"name__v": "multi-action-demo", field: "initial"}},
    )
    rec = body.get("record_id") or body.get("id")
    if not rec:
        die(f"create missing record_id: {body}")

    want_name_suffix = os.environ.get("DEMO_EXPECT_NAME_SUFFIX", "").strip()
    if want_name_suffix:
        data = client.vivarcus("object", "get", obj, rec).get("data") or {}
        got_name = data.get("name__v")
        print(f"  after create: name__v={got_name!r}")
        if not isinstance(got_name, str) or not got_name.endswith(want_name_suffix):
            die(f"trigger name__v: want suffix {want_name_suffix!r} got {got_name!r}")
        print("PASS record trigger stamped name__v")

    print(f"=== execute {set_action} ===")
    client.api(
        "POST",
        f"/api/v1/objects/{obj}/records/{rec}/actions/execute",
        {"action": set_action},
    )
    data = client.vivarcus("object", "get", obj, rec).get("data") or {}
    got = data.get(field)
    print(f"  after set: {field}={got!r}")
    if got != want_set:
        die(f"set_title: want {want_set!r} got {got!r}")

    print(f"=== execute {clear_action} ===")
    client.api(
        "POST",
        f"/api/v1/objects/{obj}/records/{rec}/actions/execute",
        {"action": clear_action},
    )
    data = client.vivarcus("object", "get", obj, rec).get("data") or {}
    got = data.get(field)
    print(f"  after clear: {field}={got!r}")
    if got != want_clear:
        die(f"clear_title: want {want_clear!r} got {got!r}")

    entry_want = os.environ.get("DEMO_ENTRY_TITLE_EXPECTED", "").strip()
    submit = os.environ.get("DEMO_SUBMIT_ACTION", "submit__c").strip()
    if entry_want:
        print("=== verify entry_action (transition → in_review) ===")
        body = client.api(
            "POST",
            f"/api/v1/objects/{obj}/records",
            {"fields": {"name__v": "multi-entry-demo", field: "initial"}},
        )
        rec_entry = body.get("record_id") or body.get("id")
        if not rec_entry:
            die(f"create missing record_id: {body}")

        want_name_suffix = os.environ.get("DEMO_EXPECT_NAME_SUFFIX", "").strip()
        if want_name_suffix:
            data = client.vivarcus("object", "get", obj, rec_entry).get("data") or {}
            got_name = data.get("name__v")
            print(f"  after create: name__v={got_name!r}")
            if not isinstance(got_name, str) or not got_name.endswith(want_name_suffix):
                die(f"entry create trigger name__v: want suffix {want_name_suffix!r} got {got_name!r}")

        data = client.vivarcus("object", "get", obj, rec_entry).get("data") or {}
        got = data.get(field)
        print(f"  before transition: {field}={got!r}")
        if got == entry_want:
            die(f"entry_action should not run before transition, got {field}={got!r}")

        client.vivarcus("lifecycle", "transition", obj, rec_entry, "--action", submit)
        data = client.vivarcus("object", "get", obj, rec_entry).get("data") or {}
        got = data.get(field)
        state = data.get("state__v")
        print(f"  after transition: {field}={got!r} state={state!r}")
        if got != entry_want:
            die(f"entry_action: want {entry_want!r} got {got!r} state={state!r}")
        print("PASS entry_action stamped title on state enter")

    print("PASS multi-component demo")


def main() -> None:
    scenario = os.environ.get("DEMO_SCENARIO", "").strip()
    if not scenario:
        die("DEMO_SCENARIO required (user-action|lifecycle|multi-component)")
    if scenario == "user-action":
        verify_user_action()
    elif scenario == "lifecycle":
        verify_lifecycle()
    elif scenario == "multi-component":
        verify_multi_component()
    else:
        die(f"unknown DEMO_SCENARIO={scenario!r}")


if __name__ == "__main__":
    main()
