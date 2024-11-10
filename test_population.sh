#!/bin/bash

mkdir test_dir
cd test_dir

touch main.go
touch main.py
touch main.c
touch Main.java
touch main.rb
touch main.rs
touch main.cs
touch main.js
touch main.ts
touch main.php
touch index.html
touch styles.css
touch script.sh
touch main.kt
touch main.swift
touch main.scala
touch main.pl
touch main.r
touch main.lua
touch Main.hs
touch main.m
touch main.groovy
touch main.dart
touch main.ex
touch main.erl
touch main.f90
touch main.pas
touch main.m
touch main.jl
touch main.sql
touch data.json
touch data.xml
touch main.sql
touch main.vhdl
touch main.cob
touch main.asm
touch main.as
touch main.vim
touch main.bash
touch main.ada
touch main.dpr
touch main.st
touch main.scm
touch main.clj
touch main.fs
touch main.ml
touch main.nim
touch main.rkt
touch main.cpp

# Go
echo -e "// This is a comment\n\npackage main\n\nfunc main() {\n}" > main.go

# Python
echo -e "# This is a comment\n\n\ndef main():\n\tpass" > main.py

# C
echo -e "// This is a comment\n\n#include <stdio.h>\n\nint main() {\n\treturn 0;\n}" > main.c

# Java
echo -e "// This is a comment\n\npublic class Main {\n\tpublic static void main(String[] args) {\n\t}\n}" > Main.java

# Ruby
echo -e "# This is a comment\n\n\ndef main\nend" > main.rb

# Rust
echo -e "// This is a comment\n\nfn main() {\n}" > main.rs

# C#
echo -e "// This is a comment\n\nclass Program {\n\tstatic void Main() {\n\t}\n}" > main.cs

# JavaScript
echo -e "// This is a comment\n\nfunction main() {\n}" > main.js

# TypeScript
echo -e "// This is a comment\n\nfunction main(): void {\n}" > main.ts

# PHP
echo -e "<?php\n// This is a comment\n\nfunction main() {\n}\n?>" > main.php

# HTML
echo -e "<!-- This is a comment -->\n\n<!DOCTYPE html>\n<html>\n<head>\n\t<title>Test</title>\n</head>\n<body>\n</body>\n</html>" > index.html

# CSS
echo -e "/* This is a comment */\n\nbody {\n}" > styles.css

# Shell
echo -e "# This is a comment\n\n\nmain() {\n}" > script.sh

# Kotlin
echo -e "// This is a comment\n\nfun main() {\n}" > main.kt

# Swift
echo -e "// This is a comment\n\nfunc main() {\n}" > main.swift

# Scala
echo -e "// This is a comment\n\nobject Main extends App {\n}" > main.scala

# Perl
echo -e "# This is a comment\n\n\nsub main {\n}" > main.pl

# R
echo -e "# This is a comment\n\n\nmain <- function() {\n}" > main.r

# Lua
echo -e "-- This is a comment\n\n\nfunction main()\nend" > main.lua

# Haskell
echo -e "-- This is a comment\n\n\nmain = do\n" > Main.hs

# Objective-C
echo -e "// This is a comment\n\n#import <Foundation/Foundation.h>\n\nint main(int argc, const char * argv[]) {\n\t@autoreleasepool {\n\t}\n\treturn 0;\n}" > main.m

# Groovy
echo -e "// This is a comment\n\nclass Main {\n\tstatic void main(String[] args) {\n\t}\n}" > main.groovy

# Dart
echo -e "// This is a comment\n\nvoid main() {\n}" > main.dart

# Elixir
echo -e "# This is a comment\n\n\ndefmodule Main do\n\tdef main do\n\tend\nend" > main.ex

# Erlang
echo -e "% This is a comment\n\n\n-module(main).\n-export([main/0]).\n\nmain() ->\n." > main.erl

# Fortran
echo -e "! This is a comment\n\n\nprogram main\nend program main" > main.f90

# Pascal
echo -e "// This is a comment\n\nprogram Main;\nbegin\nend." > main.pas

# MATLAB
echo -e "% This is a comment\n\n\nfunction main()\nend" > main.m

# Julia
echo -e "# This is a comment\n\n\nfunction main()\nend" > main.jl

# SQL
echo -e "-- This is a comment\n\n\nSELECT 1;" > main.sql

# JSON
echo -e "{\n\t\"comment\": \"This is a comment\"\n}" > data.json

# XML
echo -e "<!-- This is a comment -->\n\n<root>\n</root>" > data.xml

# T-SQL
echo -e "-- This is a comment\n\n\nSELECT 1;" > main.sql

# VHDL
echo -e "-- This is a comment\n\n\nentity main is\nend main;\n\narchitecture Behavioral of main is\nbegin\nend Behavioral;" > main.vhdl

# COBOL
echo -e "* This is a comment\n\n\nIDENTIFICATION DIVISION.\nPROGRAM-ID. Main.\n\nPROCEDURE DIVISION.\n\tSTOP RUN." > main.cob

# Assembly
echo -e "; This is a comment\n\n\nsection .text\n\tglobal _start\n\n_start:\n" > main.asm

# ActionScript
echo -e "// This is a comment\n\n\nfunction main() {\n}" > main.as

# VimL
echo -e "\" This is a comment\n\n\nfunction! Main()\nendfunction" > main.vim

# Bash
echo -e "# This is a comment\n\n\nmain() {\n}" > main.bash

# Ada
echo -e "-- This is a comment\n\n\nprocedure Main is\nbegin\nend Main;" > main.ada

# Delphi
echo -e "// This is a comment\n\n\nprogram Main;\nbegin\nend." > main.dpr

# Smalltalk
echo -e "\" This is a comment\n\n\nmain\n" > main.st

# Scheme
echo -e "; This is a comment\n\n\n(define (main)\n)" > main.scm

# Clojure
echo -e "; This is a comment\n\n\n(defn main []\n)" > main.clj

# F#
echo -e "// This is a comment\n\n\nlet main () =\n" > main.fs

# OCaml
echo -e "(* This is a comment *)\n\n\nlet main () =\n" > main.ml

# Nim
echo -e "# This is a comment\n\n\nproc main() =\n" > main.nim

# Racket
echo -e "; This is a comment\n\n\n(define (main)\n)" > main.rkt

# C++
echo -e "// This is a comment\n\n\n#include <iostream>\n\nint main() {\n\treturn 0;\n}" > main.cpp