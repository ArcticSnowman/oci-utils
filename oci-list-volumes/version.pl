#!/usr/bin/perl
use POSIX;
use strict;
use warnings;

my $version = $ENV{VERSION} // `git rev-parse --abbrev-ref HEAD`;
chomp $version;
unless ($version) {
        die "Can't find version!\n";
}

my $name    = "oci-list-volumes";
my $build   = POSIX::strftime("%Y%m%d.%H%M%S", gmtime time);
my $hash    = `git rev-parse --short HEAD`;
chomp($hash);

# Write out the package.
my $fname = "version.go";
open my $fh, "> $fname"
    or die "Can't open $fname: $!";

print $fh <<__;
// Code generated automatically by `go generate`. DO NOT EDIT.
//
// Committing updated versions of this file to git is fine but unnecessary;
// `go generate` is run while creating the Docker image, so the checked-in
// version is never actually used in deployed environments.
package ocilistvolumes

//go:generate /usr/bin/perl version.pl
const VERSION = "$build#$hash"
const DAEMON  = `$name $version Build $build#$hash`
__
    
