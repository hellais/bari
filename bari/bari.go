package main

import (
  "os"
  "log"
  "fmt"
  "os/exec"
  "strings"
  "runtime"
  "regexp"
	"encoding/json"
  "github.com/codegangsta/cli"
  "github.com/mcuadros/go-version"
)

type Package map[string]interface{};

type Packages []Package;

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
  os, ok := pkg["os"];
  if (ok && os != "multi" && os != p.os) {
    return false;   
  }
  distro, ok := pkg["distro"];
  if (ok && distro != p.distro) {
    return false;   
  }
  release, ok := pkg["release"];
  if (ok && p.release_supported(release) == false) {
  }
  //pkg_manager, ok := pkg["pkg_manager"];
  //arch, ok := pkg["arch"];
  return true;
}

func install_package(name string) {
  var platform Platform;
  platform.detect();

  pkgsDirectory := "/Users/x/code/security/install.cat/packages/";
  pkgFilePath := pkgsDirectory + name + ".json";
  pkgFile, err := os.Open(pkgFilePath);

  if err != nil {
    log.Fatal("Error in opening descriptor for package ", name);
  }

  var pkgs Packages;
  var supportedPackages Packages;

  decoder := json.NewDecoder(pkgFile);
  err = decoder.Decode(&pkgs);
  if err != nil {
    log.Fatal(err);
    log.Fatal("Error in decoding descriptor for package ", name);
  }
  
  for _, pkg := range pkgs {
    if platform.supports(pkg) {
      supportedPackages = append(supportedPackages, pkg);
    }
  }
  
  if len(supportedPackages) == 0 {
    log.Fatal("Could not find a supported install method for your OS");
  }
  
  fmt.Println("# Install options");
  for _, pkg := range supportedPackages {
    pkg_manager, ok := pkg["pkg_manager"];
    if (ok) {
      fmt.Println("##", pkg_manager);
    }
    pkg_manager_version, ok := pkg["pkg_manager_version"];
    if (ok) {
      fmt.Println("Package manager version:", pkg_manager_version);
    }

    OS, ok := pkg["os"];
    if (ok) {
      fmt.Println("Operating system:", OS);
    }

    release, ok := pkg["release"];
    if (ok) {
      fmt.Println("Release:", release);
    }

    pkg_name, ok := pkg["pkg"];
    if (ok) {
      fmt.Println("Package name:", pkg_name);
    }

  }

}

func main() {
  app := cli.NewApp()
  app.Name = "bari"
  app.Usage = "Install everything, everywhere."
  app.Commands = []cli.Command{
    {
      Name: "install",
      Usage: "Install the specified package on your system",
      Action: func(c *cli.Context) {
        pkg_name := c.Args().First();
        if pkg_name == "" {
          cli.ShowAppHelp(c);
          log.Fatal("You MUST specify a package name");
        }
        //log.Print("Installing ", pkg_name);
        install_package(pkg_name);
      },
    },
  };
  app.Run(os.Args)
}
