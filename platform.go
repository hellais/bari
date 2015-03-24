package main

import (
  "os/exec"
  "strings"
  "runtime"
  "regexp"
  
  "github.com/mcuadros/go-version"
)


type Platform struct {
  os string;
  release string;
  arch string;
  distro string;
}

func (p *Platform) detect_distro() {
  cmd := exec.Command("lsb_release", "-a");
  output, err := cmd.CombinedOutput();
  if err != nil {
    cmd := exec.Command("cat", "/etc/redhat-release");
    _, err := cmd.CombinedOutput();
    if err == nil {
      p.distro = "fedora";
    } else {
      p.distro = "unknown";
    }
  } else {
    for _, line := range strings.Split(string(output), "\n") {
      if strings.HasPrefix(line, "Distributor ID") {
        p.distro = strings.ToLower(strings.TrimSpace(strings.Split(line, ":")[1]));
      }
      if strings.HasPrefix(line, "Release") {
        p.release = strings.ToLower(strings.TrimSpace(strings.Split(line, ":")[1]));
      }
    }
  }
}

func (p *Platform) detect_osx_release() {
  cmd := exec.Command("sw_vers", "-productVersion");
  output, _ := cmd.CombinedOutput();
  p.release = strings.TrimSpace(string(output));
}


func (p *Platform) detect() {
  switch runtime.GOOS {
    case "darwin":
      p.os = "osx";
    case "windows":
      p.os = "windows";
    case "linux":
      p.os = "linux";
    default:
      p.os = "unknown";
  }

  if p.os == "linux" {
    p.detect_distro();
  } else if p.os == "osx" {
    p.detect_osx_release();
  }
}

func (p *Platform) release_supported(release interface{}) bool {
  var result bool;
  r, _ := regexp.Compile(`([>|<][=]?)?((\d)(\.\d))+?`);
  parts := r.FindAllStringSubmatch(release.(string), -1)[0];
  target_version := parts[1];
  comparison := parts[0];
  if comparison != "" {
    result = version.Compare(p.release, target_version, comparison)
  } else {
    result = version.Compare(p.release, target_version, ">=")
  }
  return result;
}

func (p *Platform) supports(pkg Package) bool {
  if (pkg.os != "multi" && pkg.os != p.os) {
    return false;
  }
  if (pkg.distro != p.distro) {
    return false;
  }
  if (pkg.release != "" && p.release_supported(pkg.release)) {
    return false;
  }
  //pkg_manager, ok := pkg["pkg_manager"];
  //arch, ok := pkg["arch"];
  return true;
}
