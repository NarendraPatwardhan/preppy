from sphinx.ext.intersphinx import fetch_inventory
import json

URL = "https://docs.python.org/{}/objects.inv"
VERSIONS = [("3", "8"), ("3", "9"), ("3", "10")]

# Create a placeholder config for sphinx
class Config:
    intersphinx_timeout = None
    tls_verify = True
    user_agent = ""

# Create a placeholder app for sphinx
class App:
    srcdir = ""
    config = Config()

# Create a dictionary to hold core packages
modules = {}

# Iterate over the versions and fetch the inventory
for version_info in VERSIONS:
    version = ".".join(version_info)
    url = URL.format(version)
    invdata = fetch_inventory(App(), "", url)
    # Get the core packages
    for module in invdata["py:module"]:
        root, *_ = module.split(".")
        # If not already accounted for, add it to the dictionary
        if root not in modules:
            modules[root] = True

# Write modules as a json file
with open("core/stdlib.json", "w") as f:
    json.dump(modules, f)