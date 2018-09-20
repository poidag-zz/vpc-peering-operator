#!/usr/bin/make -f

make: build

build:
		operator-sdk generate k8s
		operator-sdk build quay.io/pickledrick/vpc-peering-operator