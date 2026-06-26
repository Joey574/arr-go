package fsu

import (
	"arr-go/v2/internal/log"
	"io"
	"os"
)

func Symlink(src, dst string) error {
	srcinfo, srcerr := os.Stat(src)
	dstinfo, dsterr := os.Stat(dst)

	if srcerr == nil && srcinfo.IsDir() {
		return log.AsError("src is a directory")
	}

	if dsterr == nil && dstinfo.IsDir() {
		return log.AsError("dst is a directory")
	}

	if srcerr != nil && dsterr != nil {
		// errors opening src and dst, cannot recover
		return log.AsError("src: %v, dst: %v, cannot recover, exiting", srcerr, dsterr)
	} else if srcerr != nil && dsterr == nil {
		// src doesn't exist, dst does, recoverable
		log.Warnf("src: %v, dst: found, recoverable", srcerr)
		err := os.Chmod(dst, 0770)
		if err != nil {
			log.Warnf("failed to change dst to 0770: %v", err)
		}

	} else if srcerr == nil && dsterr != nil {
		// dst doesn't exist, src does, recoverable
		log.Warnf("src: found, dst: %v, recoverable", dsterr)

		srcf, err := os.Open(src)
		if err != nil {
			return log.AsError("failed to open src for reading: %v", err)
		}

		dstf, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0770)
		if err != nil {
			return log.AsError("failed to open dst for writing: %v", err)
		}

		_, err = io.Copy(dstf, srcf)
		if err != nil {
			return log.AsError("failed to copy from src to dst: %v", err)
		}
		log.Infof("copied from src to dst")

		err = os.Remove(src)
		if err != nil {
			log.Warnf("failed to remove src: %v", err)
		} else {
			log.Infof("removed src, ready to create symlink")
		}

	} else if srcerr == nil && dsterr == nil {
		// both dst and src exist, this is the expected state
		log.Infof("src: found, dst: found")

		err := os.Remove(src)
		if err != nil {
			log.Warnf("failed to remove src: %v", err)
		} else {
			log.Infof("removed src, ready to create symlink")
		}
	}

	// remove any hanging symlinks from previous runs
	_ = os.Remove(src)

	// we have either acheived expected state or failed by now ( src doesn't exist, dst does )
	err := os.Symlink(dst, src)
	if err != nil {
		return log.AsError("failed to create symlink from src -> dst: %v", err)
	}

	log.Infof("created symlink from src -> dst")
	return nil
}
