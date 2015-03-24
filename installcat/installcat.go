package main

import (
  "os"
  "log"
  "fmt"
	"encoding/json"
  "github.com/codegangsta/cli"
)

func install_package(name string) {
  pkgsDirectory := "/Users/x/code/security/install.cat/packages/autogen/";
  pkgFilePath := pkgsDirectory + name + ".json";
  pkgFile, err := os.Open(pkgFilePath);

  if err != nil {
    log.Fatal("Error in opening descriptor for package", name);
  }

  type Packages []map[string]string;

  var pkgs Packages;
  decoder := json.NewDecoder(pkgFile);
  err = decoder.Decode(&pkgs);
  if err != nil {
    log.Fatal(err);
    log.Fatal("Error in decoding descriptor for package", name);
  }
  
  for idx := range pkgs {
    pkg := pkgs[idx];
    for key, value := range pkg {
        fmt.Println(key, ":", value)
    }
    fmt.Println("--------");
  }
}

func main() {
  app := cli.NewApp()
  app.Name = "install"
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
        log.Print("Installing", pkg_name);
        install_package(pkg_name);
      },
    },
  };
  app.Run(os.Args)
}
