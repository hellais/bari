package main

import (
  "os"
  "log"
  "fmt"
	"encoding/json"

  "github.com/codegangsta/cli"
)

type Package map[string]interface{};

type Packages []Package;


func install_package(name string) {
  var platform Platform;
  platform.detect();

  pkgsDirectory := "packages/";
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
      Aliases: []string{"i"},
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

    {
      Name: "show",
      Usage: "Show the setup instructions for the specified package",
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
