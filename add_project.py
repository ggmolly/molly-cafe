#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import json
import argparse
import os

PROJECT_FOLDER = os.path.join("projects", "school")
os.makedirs(PROJECT_FOLDER, exist_ok=True)

TEMPLATE = {
    "name": "",
    "href": "",
    "description": "",
    "wip": False,
    "grading": False,
    "grade": 0,
}

if __name__ == "__main__":
    args = argparse.ArgumentParser()
    args.add_argument("-n", "--name", help="Name of the project", required=True)
    args.add_argument("-d", "--description", help="Description of the project", default="")
    args.add_argument("-l", "--link", help="Link to the project (GitHub / Pistache post)", default="")
    args.add_argument("-w", "--wip", help="Work in progress", action="store_true", default=False)
    args.add_argument("-g", "--grading", help="Grading", action="store_true", default=False)
    args.add_argument("-gr", "--grade", help="Grade", type=int, default=0)
    args = args.parse_args()

    template = TEMPLATE.copy()
    template["name"] = args.name
    template["description"] = args.description
    template["href"] = args.link
    template["wip"] = args.wip
    template["grading"] = args.grading
    template["grade"] = args.grade

    with open(os.path.join(PROJECT_FOLDER, args.name + ".json"), "w+") as f:
        json.dump(template, f, indent=4, sort_keys=True)

    print("Project added successfully!")