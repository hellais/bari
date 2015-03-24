package main

import (
  "os"
  "log"

	"encoding/json"
)

type PackageJSON map[string]interface{};

type Package struct {
  os string;
  distro string;
  release string;
  arch []string;
  pkg_manager string;
  pkg_manager_version string;
  pkg string;
  priority int;
  repo string;
  url string;
  json PackageJSON;
}

func (p *Package) install_command() []string {
  var command []string;
  if (p.pkg_manager == "apt") {
    command = []string{"apt-get", "install", p.pkg};
  } else if (p.pkg_manager == "yum") {
    command = []string{"yum", "install", p.pkg};
  } else if (p.pkg_manager == "homebrew") {
    command = []string{"brew", "install", p.pkg};
  } else if (p.pkg_manager == "pip") {
    command = []string{"pip", "install", p.pkg};
  } else if (p.pkg_manager == "pacman") {
    command = []string{"pacman", "install", p.pkg};
  }
  return command;
}

func NewPackage(packageJSON PackageJSON) Package {
  p := new(Package);
  p.json = packageJSON;

  OS, ok := packageJSON["os"];
  if (ok) {
    p.os = OS.(string);
  } else {
    p.os = "";
  }

  release, ok := packageJSON["release"];
  if (ok) {
    p.release = release.(string);
  } else {
    p.release = "";
  }

  distro, ok := packageJSON["distro"];
  if (ok) {
    p.distro = distro.(string);
  } else {
    p.distro = "";
  }

  pkg_manager, ok := packageJSON["pkg_manager"];
  if (ok) {
    p.pkg_manager = pkg_manager.(string);
  } else if (p.os == "linux") {
    if (p.distro == "debian" ||
        p.distro == "ubuntu") {
      p.pkg_manager = "apt";
    } else if (p.distro == "archlinux") {
      p.pkg_manager = "pacman";
    } else if (p.distro == "centos" ||
               p.distro == "fedora" ||
               p.distro == "redhat") {
      p.pkg_manager = "yum";
    }
  } else if (p.os == "osx") {
    p.pkg_manager = "homebrew";
  }

  pkg_manager_version, ok := packageJSON["pkg_manager_version"];
  if (ok) {
    p.pkg_manager_version = pkg_manager_version.(string);
  } else {
    p.pkg_manager_version = "";
  }

  pkg_name, ok := packageJSON["pkg"];
  if (ok) {
    p.pkg = pkg_name.(string);
  } else {
    p.pkg = "";
  }

  return *p;
}

func loadPackages(name string) []Package {
  pkgsDirectory := "packages/";
  pkgFilePath := pkgsDirectory + name + ".json";
  pkgFile, err := os.Open(pkgFilePath);

  if err != nil {
    log.Fatal("Error in opening descriptor for package ", name);
  }

  var packagesJSON []PackageJSON;
  var packages []Package;

  decoder := json.NewDecoder(pkgFile);
  err = decoder.Decode(&packagesJSON);

  if err != nil {
    log.Fatal(err);
    log.Fatal("Error in decoding descriptor for package ", name);
  }

  for _, packageJSON := range packagesJSON {
    packages = append(packages, NewPackage(packageJSON));
  }
  return packages;
}
