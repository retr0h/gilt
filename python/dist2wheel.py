#!/usr/bin/env python3

# Assemble Python wheels for all artifacts built under dist/, according to the
# PEP 491 specification.  Does the absolute bare minimum to get
# `pip install` and `pip uninstall` to place the gilt binary under bin/.
#
# Extraordinarily opinionated about where things live, and not
# expected to be usable by anything other than Gilt.

# These packages should be included in the Python stdlib.  Nothing
# other than a basic Python interpreter should be needed.
import base64
import hashlib
import json
import os
import zipfile


# The METADATA file can likely be made even smaller, but explicit
# is better than implicit.  ;-)
WHEEL_TEMPLATE = """Wheel-Version: 1.0
Generator: dist2wheel.py
Root-Is-Purelib: false
Tag: {py_tag}-{abi_tag}-{platform_tag} 
"""
METADATA_TEMPLATE = """Metadata-Version: 2.1
Name: {distribution}
Version: {version}
Classifier: Development Status :: 5 - Production/Stable
Classifier: Environment :: Console
Classifier: Intended Audience :: Developers
Classifier: Operating System :: OS Independent
Classifier: Programming Language :: Python
Classifier: Programming Language :: Python :: 3 :: Only
Classifier: Programming Language :: Go
Summary: {description}
License: {license}
Description-Content-Type: text/markdown; charset=UTF-8; variant=GFM

{readme}
"""


class Wheels:
    def __init__(self):
        # Pull in all the project metadata from known locations.
        # No, this isn't customizable.  Cope.
        with open("dist/metadata.json") as f:
            self.metadata = json.load(f)

        with open("dist/artifacts.json") as f:
            self.artifacts = json.load(f)

        with open("README.md") as f:
            self.readme = f.read()

        self.distribution = self.metadata["project_name"]
        self.version = self.metadata["version"]
        self.py_tag = "py3"
        self.abi_tag = "none"

    def create_all(self):
        """Generate a Python wheel for every artifact in artifacts.json."""

        # Burn a pass through the list to steal some useful bits from the Brew config
        for artifact in self.artifacts:
            if "BrewConfig" in artifact["extra"]:
                self.description = artifact["extra"]["BrewConfig"]["description"]
                self.license = artifact["extra"]["BrewConfig"]["license"]

        # We're looking for "internal_type: 2" artifacts, but being an internal
        # type, we'll avoid leaning on implementation details if we don't have to
        for artifact in self.artifacts:
            try:
                self.path = artifact["path"]
                self.platform_tag = self._fixup_platform_tag(artifact)
                self.checksum = self._fix_checksum(artifact["extra"]["Checksum"])
                self.size = os.path.getsize(self.path)
            except KeyError:
                continue

            self.wheel_file = WHEEL_TEMPLATE.format(**self.__dict__).encode()
            self.metadata_file = METADATA_TEMPLATE.format(**self.__dict__).encode()
            self._emit()

    @staticmethod
    def _fixup_platform_tag(artifact):
        """Convert Go binary nomenclature to Python wheel nomenclature."""

        # Go 1.21 will require macOS 10.15 or later
        _map = dict(darwin="macosx_10_15", linux="linux")
        platform = _map[artifact["goos"]]

        arch = artifact["goarch"]
        if arch == "arm64" and platform == "linux":
            arch = "aarch64"
        elif arch == "amd64":
            arch = "x86_64"
        elif arch == "386":
            arch = "i686"

        return f"{platform}_{arch}"

    @staticmethod
    def _fix_checksum(checksum):
        """Re-encode the checksum as base64, with no trailing = characters."""

        if checksum.startswith("sha256:"):
            checksum = checksum[7:]
        return base64.urlsafe_b64encode(bytes.fromhex(checksum)).decode().rstrip("=")

    def _emit(self):
        name_ver = f"{self.distribution}-{self.version}"
        filename = f"dist/{name_ver}-{self.py_tag}-{self.abi_tag}-{self.platform_tag}.whl"
        with zipfile.ZipFile(filename, "w", compression=zipfile.ZIP_DEFLATED) as zf:
            record = []
            print(f"writing {zf.filename}")

            # The actual binary on-disk, and the recorded checksum from artifacts.json
            arcname = f"{name_ver}.data/scripts/{os.path.basename(self.path)}"
            zf.write(self.path, arcname=arcname)
            record.append(f"{arcname},sha256={self.checksum},{self.size}")

            # The project metadata
            arcname = f"{name_ver}.dist-info/METADATA"
            zf.writestr(arcname, self.metadata_file)
            digest = hashlib.sha256(self.metadata_file).hexdigest()
            record.append(f"{arcname},sha256={self._fix_checksum(digest)},{len(self.metadata_file)}")

            # The platform tags
            arcname = f"{name_ver}.dist-info/WHEEL"
            zf.writestr(arcname, self.wheel_file)
            digest = hashlib.sha256(self.wheel_file).hexdigest()
            record.append(f"{arcname},sha256={self._fix_checksum(digest)},{len(self.wheel_file)}")

            # Write out the manifest last.  The record of itself contains no checksum or size info
            arcname = f"{name_ver}.dist-info/RECORD"
            record.append(f"{arcname},,")
            zf.writestr(arcname, "\n".join(record))


if __name__ == "__main__":
    Wheels().create_all()
