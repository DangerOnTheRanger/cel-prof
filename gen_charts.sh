#!/bin/sh
go tool pprof -alloc_space -png mem.prof
go tool pprof -png cpu.prof
