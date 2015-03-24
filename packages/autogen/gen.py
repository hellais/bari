import re
import sys
import json
import time

from pprint import pprint
from distutils.version import LooseVersion

import requests
from bs4 import BeautifulSoup

# # soup.select(".ui-accordion-content")
# distro_version_pairs = []
# full_distro_names = [distro.text for distro in soup.select(".distro")]
# for distro in full_distro_names:
#     m = re.search("([^\d]+)((\s+)?\d.*)?", distro).groups()
#     version = None
#     if m[1] is not None:
#         version = "%s%s" % (m[2], m[1]) if m[2] is not None else m[1]
#     name = m[0]
#     # if mappings.get(name) is None:
#     #     mappings[name] = {
#     #         "name": "XXX",
#     #         "versions": {}
#     #     }
#     # mappings[name]["versions"][version] = "XXX"
#     distro_version_pairs.append((name, version))

mappings = {
    u'ALT Linux Sisyphus': {'name': 'alt', 'versions': {None: 'XXX'}},
    u'Arch Linux': {'name': 'arch', 'versions': {None: 'XXX'}},
    u'CentOS ': {'name': 'centos', 'versions': {u'6': "6", u'7': "7"}},
    u'Debian Jessie': {'name': 'debian','version': "8.0", 'versions': {None: 'XXX'}},

    u'Debian Sid': {'name': 'debian', 'version': 'dev', 'versions': {None: 'XXX'}},
    u'Debian Wheezy': {'name': 'debian', 'version': '7.0', 'versions': {None: 'XXX'}},
    u'Fedora ': {'name': 'fedora',
                 'versions': {u'20': '20', u'21': '21', u'22': '22'}},
    u'Fedora Rawhide': {'name': 'fedora', 'version': 'dev', 'versions': {None: 'XXX'}},
    u'OpenMandriva Cooker': {'name': 'mandriva', 'version': 'dev', 'versions': {None: 'XXX'}},
    u'OpenMandriva Lx ': {'name': 'mandriva',
                          'versions': {u'2013.0': '2013.0', u'2014.1': '2014.1'}},
    u'ROSA ': {'name': 'rosa', 'versions': {u'2014.1': '2014.1'}},
    u'ROSA Desktop Fresh R': {'name': 'rosa', 'versions': {u'3': '3'}},
    u'Slackware ': {
        'name': 'slackware',
        'versions': {u'13.37': '13.37',
                     u'14.0': '14.0',
                     u'14.1': '14.1'}
    },
    u'Ubuntu ': {'name': 'ubuntu',
                 'versions': {u'12.04 LTS': '12.04',
                              u'14.04 LTS': '14.04',
                              u'14.10': '14.10'}},
    u'openSUSE ': {'name': 'suse', 'versions': {u'13.1': '13.1', u'13.2': '13.2'}},
    u'openSUSE Factory': {'name': 'suse', 'version': 'dev', 'versions': {None: 'XXX'}}}


def pkg_names_for_distro(distro_name, version, pkgs_repo):
    pkg_repo_mapping = {}
    names = []

    for pkg, repo in pkgs_repo:
        m = re.search("(.+?)((\d)(\.\d)+)", pkg)
        if not m:
            continue
        groups = m.groups()
        pkg_name = re.sub("[-_]$", "", groups[0])
        pkg_version = groups[1]
        if not pkg_repo_mapping.get("pkg_name"):
            pkg_repo_mapping[pkg_name] = []
        pkg_repo_mapping[pkg_name].append((pkg_version, repo, pkg))

    distro = mappings[distro_name]
    distro_name = distro["name"]
    if version is not None:
        try:
            distro_version = distro["versions"][version]
        except:
            distro_version = version.strip()
    elif distro.get("version"):
        distro_version = distro.get("version")
    else:
        distro_version = version

    for pkg_name, versions in pkg_repo_mapping.items():
        best_version = versions[0][0]
        best_repo = versions[0][1]
        best_full_name = versions[0][2]
        for version in versions:
            if LooseVersion(str(version)) > LooseVersion(str(best_version)):
                best_version = version[0]
                best_repo = version[1]
                best_full_name = version[2]
        name = {
            "os": "linux",
            "distro": distro_name,
            "pkg": pkg_name,
            "pkg_version": best_version,
            "pkg_repo": best_repo,
            "pkg_full_name": best_full_name
        }
        if distro_version is not None:
            name["release"] = ">=%s" % distro_version
        names.append(name)
    return names


def pkg_names(name):
    names = []

    r = requests.get("http://pkgs.org/download/%s" % name)
    soup = BeautifulSoup(r.text)

    pkgs_repo = []
    distro_version = None
    distro_name = None
    for div in soup.select("#pkgs_show div"):
        if "distro" in div.get("class"):
            if distro_name is not None:
                if distro_name in mappings:
                    names += pkg_names_for_distro(distro_name, distro_version, pkgs_repo)
            distro_name = div.text
            m = re.search("([^\d]+)((\s+)?\d.*)?", distro_name).groups()
            distro_name = m[0]
            distro_version = None
            if m[1] is not None:
                distro_version = "%s%s" % (m[2], m[1]) if m[2] is not None else m[1]
        elif "ui-accordion-content" in div.get("class"):
            repo = None
            for sub_div in div.select("*"):
                if not sub_div.get("class"):
                    continue
                if "repo" in sub_div.get("class"):
                    repo = sub_div.text.replace(":", "")
                elif "level2" in sub_div.get("class"):
                    for pkg in sub_div.select("li a"):
                        pkgs_repo.append((pkg.text, repo))
    return names

def get_global_name(url):
    r = requests.get(url)
    soup = BeautifulSoup(r.text)
    return soup.select("#see_also li a")[0].text

def get_packages_in_page(path):
    r = requests.get("http://pkgs.org%s" % path)
    soup = BeautifulSoup(r.text)
    return [x.get("href") for x in soup.select("#pkgs_show li a")]

def list_all_packages(url="http://pkgs.org/debian-sid/debian-main-i386/"):
    r = requests.get(url)
    soup = BeautifulSoup(r.text)

    paths = [x.get("href") for x in soup.select("#pkgs_show li a")]
    # resumed = False
    # for path in paths:
    #     if resumed is False and "alsa-base" not in path:
    #         continue
    #     resumed = True
    #     name = get_global_name("http://pkgs.org%s" % path)
    #     print "Operating on %s" % name
    #     with open("%s.json" % name, "w+") as f:
    #         json.dump(pkg_names(name), f, indent=2)

    pages = [x.get("href") for x in soup.select("#pkgs_show-pages a")]
    resumed = False
    for page in pages:
        print "Getting page %s" % page
        if "/3/" not in page and resumed is False:
            continue
        resumed = True
        paths = get_packages_in_page(page)
        for path in paths:
            time.sleep(1)
            name = get_global_name("http://pkgs.org%s" % path)
            print "Operating on %s" % name
            with open("%s.json" % name, "w+") as f:
                json.dump(pkg_names(name), f, indent=2)

# print get_global_name("http://pkgs.org/debian-sid/debian-main-i386/0ad_0.0.18-1_i386.deb.html")
list_all_packages()
# pprint(pkg_names(sys.argv[1]))
