#!/usr/bin/env python3
"""Generate the GL.iNet firmware 4.x API reference under docs/api/.

GL.iNet's official 4.x API reference (dev.gl-inet.com/router-4.x-api) is no longer
publicly reachable. The machine-readable description survives inside the
python-glinet package, which exports GL.iNet's own documentation database.

This script reads that description and emits our own Markdown reference. The
source file is deliberately NOT vendored into this repository: python-glinet is
GPL-3.0 while gogl is MIT, and mixing them would be a licensing problem. API
method names and signatures are functional interface facts rather than creative
work, so documenting them here is fine; redistributing their file is a different
question, and not one worth guessing at.

Fetch the source first:

    scripts/fetch-api-description.sh

then:

    python3 scripts/generate-api-docs.py /tmp/gl-api-description.json
"""

import json
import pathlib
import re
import sys
from collections import OrderedDict

REPO = pathlib.Path(__file__).resolve().parent.parent
OUT = REPO / "docs" / "api"

# Endpoints confirmed by calling them against a live GL-SFT1200 running firmware
# 4.3.28 on 2026-07-28. Everything else in the reference is from the description
# and has not been exercised here.
VERIFIED = {
    ("system", "get_info"),
    ("system", "get_status"),
    ("clients", "get_list"),
    ("clients", "get_status"),
    ("lan", "get_config_list"),
    ("lan", "get_static_bind_list"),
    ("lan", "add_static_bind"),
    ("lan", "set_static_bind"),
    ("lan", "remove_static_bind"),
    ("network", "get_dhcp_leases"),
    ("wifi", "get_config"),
    ("wifi", "get_status"),
    ("wifi", "set_config"),
    ("dns", "get_config"),
    ("dns", "get_host"),
    ("dns", "set_host"),
    ("lan", "set_config"),
    ("repeater", "get_config"),
    ("led", "get_config"),
    ("qos", "get_config"),
    ("igmp", "get_config"),
    ("ddns", "get_config"),
    ("cloud", "get_config"),
    ("tor", "get_config"),
    ("upgrade", "get_config"),
}

# Endpoints that returned -32601 Method not found on that device, i.e. documented
# but absent from this model or firmware.
ABSENT = {
    ("custom_dns", "get_info"),
    ("custom_dns", "set_info"),
    ("modem", "get_config"),
    ("system", "get_config"),
    ("lan", "get_config"),
    ("network", "get_config"),
    ("network", "get_status"),
    ("firewall", "get_config"),
    ("firewall", "get_status"),
}


def first(value):
    """Unwrap the single-element lists the description uses for scalars."""
    if isinstance(value, list):
        return value[0] if value else ""
    return value or ""


def slug(name):
    """GitHub's heading anchor: lowercased, punctuation dropped, spaces hyphenated.

    Underscores are preserved, which matters here because every method name has
    them -- hyphenating produced links that silently went nowhere.
    """
    anchor = name.lower()
    anchor = re.sub(r"[^\w\s-]", "", anchor)
    return re.sub(r"\s+", "-", anchor).strip("-")


def pretty(raw):
    """Pretty-print a JSON example, leaving it verbatim if it will not parse."""
    if not raw:
        return None
    try:
        return json.dumps(json.loads(raw), indent=2)
    except (ValueError, TypeError):
        return raw.strip()


def status_note(group, method):
    if (group, method) in VERIFIED:
        return "**Verified** on a GL-SFT1200 running 4.3.28."
    if (group, method) in ABSENT:
        return "**Absent** on a GL-SFT1200 running 4.3.28: returns `-32601 Method not found`."
    return "Not exercised against hardware here."


def field_table(entries, kind):
    """Render params or results as a table.

    A leading '?' on a key name marks it optional in the source's convention.
    """
    if not entries:
        return [f"_No {kind}._", ""]

    lines = [f"| {kind.capitalize()} | Type | Required | Description |",
             "|---|---|---|---|"]
    for e in entries:
        key = str(e.get("keyName") or "")
        optional = key.startswith("?")
        key = key.lstrip("?")
        required = "no" if optional else "yes"
        # results have no meaningful required/optional sense
        if kind == "results":
            required = "-"
        dtype = e.get("dataType__name") or ""
        desp = (e.get("desp") or "").replace("\n", " ").replace("|", "\\|").strip()
        lines.append(f"| `{key}` | {dtype} | {required} | {desp} |")
    lines.append("")
    return lines


def render_group(group, payload):
    methods = OrderedDict(sorted((payload.get("case_groups_data") or {}).items()))
    desp = first(payload.get("module_desp"))

    lines = [
        f"# `{group}`",
        "",
        desp or "_No description in the source._",
        "",
        f"{len(methods)} method(s). Call shape:",
        "",
        "```json",
        '{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","%s","<method>",{}]}' % group,
        "```",
        "",
        "| Method | Status | Description |",
        "|---|---|---|",
    ]
    for name, m in methods.items():
        d = (m.get("data") or {}).get("desp") or ""
        mark = "verified" if (group, name) in VERIFIED else (
            "absent" if (group, name) in ABSENT else "-")
        lines.append(f"| [`{name}`](#{slug(name)}) | {mark} | {d.replace('|', chr(92) + '|')} |")
    lines.append("")
    lines.append("---")
    lines.append("")

    for name, m in methods.items():
        data = m.get("data") or {}
        lines += [
            f"## {name}",
            "",
            data.get("desp") or "_No description._",
            "",
            status_note(group, name),
            "",
        ]
        lines += field_table(m.get("params"), "params")
        lines += field_table(m.get("results"), "results")

        for label, key in (("Request", "in_example"), ("Response", "out_example")):
            example = pretty(m.get(key))
            if example:
                lines += [f"**{label}**", "", "```json", example, "```", ""]
        lines += ["---", ""]

    return "\n".join(lines)


def render_index(desc):
    groups = OrderedDict(sorted(desc.items()))
    total = sum(len(g.get("case_groups_data") or {}) for g in groups.values())
    verified_here = sum(
        1 for g, payload in groups.items()
        for m in (payload.get("case_groups_data") or {})
        if (g, m) in VERIFIED
    )

    lines = [
        "# GL.iNet Firmware 4.x JSON-RPC API Reference",
        "",
        f"{len(groups)} groups, {total} methods.",
        "",
        "## Provenance",
        "",
        "GL.iNet's official 4.x API reference at `dev.gl-inet.com/router-4.x-api` is no",
        "longer publicly reachable. This reference was generated from the",
        "machine-readable API description that ships inside",
        "[`python-glinet`](https://github.com/tomtana/python-glinet), which exports",
        "GL.iNet's own documentation database.",
        "",
        "That source file is **not** vendored here. `python-glinet` is GPL-3.0 and gogl is",
        "MIT, so redistributing it would be a licensing question not worth guessing at.",
        "Method names and signatures are functional interface facts and are documented",
        "freely below; to obtain the source yourself:",
        "",
        "```bash",
        "scripts/fetch-api-description.sh",
        "python3 scripts/generate-api-docs.py /tmp/gl-api-description.json",
        "```",
        "",
        "## What is actually verified",
        "",
        "The description is GL.iNet's documentation for the firmware line as a whole, not",
        f"for any one device. {verified_here} endpoints were confirmed by calling them",
        "against a **GL-SFT1200 (Opal) running firmware 4.3.28**, and several documented",
        "endpoints are absent on it. Each method below is marked accordingly.",
        "",
        "Two payloads differ from the documentation on real hardware:",
        "",
        "- `network.get_dhcp_leases` wraps its array in `leases`, not the documented `entries`.",
        "- `lan.get_config_list` returns more fields than documented (`leasetime`, `dns`,",
        "  `gateway`, `lpr`) and reports `enable` as a **number**, not a boolean.",
        "",
        "Treat this reference as a map, and the device as the authority. See",
        "[`../../GL_INET_4X_API_DOCUMENTATION.md`](../../GL_INET_4X_API_DOCUMENTATION.md)",
        "for authentication, error codes, and the hardware-verified essentials.",
        "",
        "## Groups",
        "",
        "| Group | Methods | Verified here | Description |",
        "|---|---|---|---|",
    ]
    for g, payload in groups.items():
        methods = payload.get("case_groups_data") or {}
        v = sum(1 for m in methods if (g, m) in VERIFIED)
        desp = str(first(payload.get("module_desp"))).replace("|", "\\|")
        lines.append(f"| [`{g}`](reference/{g}.md) | {len(methods)} | {v or '-'} | {desp} |")
    lines.append("")
    return "\n".join(lines)


def main():
    if len(sys.argv) != 2:
        print(__doc__)
        return 1

    desc = json.load(open(sys.argv[1]))
    ref = OUT / "reference"
    ref.mkdir(parents=True, exist_ok=True)

    (OUT / "README.md").write_text(render_index(desc))
    for group, payload in sorted(desc.items()):
        (ref / f"{group}.md").write_text(render_group(group, payload))

    total = sum(len(g.get("case_groups_data") or {}) for g in desc.values())
    print(f"wrote {OUT/'README.md'}")
    print(f"wrote {len(desc)} group files under {ref}")
    print(f"{total} methods documented")
    return 0


if __name__ == "__main__":
    sys.exit(main())
