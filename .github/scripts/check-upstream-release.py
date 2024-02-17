#!/usr/bin/env python3

# Check for upstream releases of flux-sched, and compare
# to current version here. If a new version is found, we

import argparse
import requests
import sys
import os

# python .github/scripts/check-upstream-release.py flux-framework/flux-sched

token = os.environ.get("GITHUB_TOKEN")
headers = {}
if token:
    headers["Authorization"] = "token %s" % token


def write_file(data, filename):
    """
    Write content to file
    """
    with open(filename, "w") as fd:
        fd.writelines(data)


def read_file(filename):
    """
    Read content from file
    """
    with open(filename, "r") as fd:
        content = fd.read()
    return content


def set_env_and_output(name, value):
    """
    helper function to echo a key/value pair to output and env.

    Parameters:
    name (str)  : the name of the environment variable
    value (str) : the value to write to file
    """
    for env_var in ("GITHUB_ENV", "GITHUB_OUTPUT"):
        environment_file_path = os.environ.get(env_var)
        print("Writing %s=%s to %s" % (name, value, env_var))

        with open(environment_file_path, "a") as environment_file:
            environment_file.write("%s=%s\n" % (name, value))


class UpstreamUpdater:
    def __init__(self, version_file, repo, dry_run=False):
        self.version_file = os.path.abspath(version_file)
        self.repo = repo
        self.dry_run = dry_run
        self._latest_version = None
        self._current_version = self.get_current_version()

    @property
    def current_version(self):
        return self._current_version

    def get_current_version(self):
        """
        Derive current version from file (or VERSION)
        """
        if not os.path.exists(self.version_file):
            sys.exit(f"{self.version_file} does not exist.")
        self._current_version = read_file(self.version_file).strip("\n")

    def check(self):
        """
        Given a repository name, check for new releases.
        """
        latest = self.get_latest_release()
        version = self.current_version
        tag = latest["tag_name"]

        # Some versions are prefixed with v
        if tag == version or tag == f"v{version}":
            print("No new version found.")
            return
        print(f"New version {tag} detected!")
        naked_version = tag.replace("v", "")
        self.update_version_file(naked_version)
        set_env_and_output("version", tag)

    def get_latest_release(self):
        """
        Get the lateset release of a repository (under flux-framework)
        """
        url = f"https://api.github.com/repos/{self.repo}/releases"
        response = requests.get(url, headers=headers, params={"per_page": 100})
        response.raise_for_status()

        # latest release should be first
        return response.json()[0]

    def update_version_file(self, version):
        """
        Update the package file with a new version and digest.
        """
        write_file(version, self.version_file)


def get_parser():
    parser = argparse.ArgumentParser(
        description="Upstream Release Updater",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--version-file",
        help="version file to parse",
        default="VERSION",
        dest="version_file",
    )
    parser.add_argument(
        "--repo", help="GitHub repository name", default="flux-framework/flux-sched"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        default=False,
        help="Don't write changes to file",
    )
    return parser


def main():
    parser = get_parser()

    # If an error occurs while parsing the arguments, the interpreter will exit with value 2
    args, extra = parser.parse_known_args()

    # Show args to the user
    print("version file: %s" % args.version_file)
    print("        repo: %s" % args.repo)
    print("     dry-run: %s" % args.dry_run)

    updater = UpstreamUpdater(args.version_file, args.repo, args.dry_run)
    updater.check()


if __name__ == "__main__":
    main()
