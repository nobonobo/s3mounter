package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func main() {
	region, bucket, rootKey := "ap-northeast-3", "nobonobo-bucket", "root"
	flag.StringVar(&region, "region", region, "AWS region")
	flag.StringVar(&bucket, "bucket", bucket, "AWS S3 bucket")
	flag.StringVar(&rootKey, "root", rootKey, "AWS S3 root key")
	flag.Parse()
	if flag.NArg() < 1 {
		log.Println("need mount point arg")
		flag.Usage()
		os.Exit(1)
	}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("LoadDefaultConfig failed: %v", err)
	}
	client := s3.NewFromConfig(cfg)
	root := &s3node{
		client: client,
		bucket: bucket,
		key:    rootKey + "/",
	}
	server, err := fs.Mount(flag.Arg(0), root, &fs.Options{
		GID: 100,
		RootStableAttr: &fs.StableAttr{
			Mode: fuse.S_IFDIR | 0775,
			Ino:  StringToIno(root.key),
		},
		MountOptions: fuse.MountOptions{
			AllowOther:     true,
			RememberInodes: false,
			DisableXAttrs:  true,
			Debug:          false,
			FsName:         "s3mount",
			Name:           "s3",
			Options:        []string{"default_permissions"},
		},
	})
	if err != nil {
		log.Fatalf("Mount failed: %v", err)
	}
	defer server.Unmount()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		log.Println("Received signal, exiting.")
		server.Unmount()
	}()
	server.Wait()
}
