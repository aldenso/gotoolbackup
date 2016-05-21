gotoolbackup
============

Program to create backups using toml file, where you indicate origin and destiny directories and retention period in days.

TODO: include log package to create and output file.

Note: this is not to replace shell scripts or anything else, is just to practice some golang and learn how to use tar and gzip packages, backups times are not bad but it could be better.

  > ./gotoolbackup --help
  Usage of ./gotoolbackup:
    -keepfiles
          indicate if you want to keep original files. (default true)
    -tomlfile string
          indicate tomlfile to read backups details. (default "gobackup.toml")


  > ./gotoolbackup --tomlfile=/tmp/gobackup.toml --keepfiles=false
  #### Running with values ####
  tomlfile: /tmp/gobackup.toml
  keepfiles: false
  #####################################
  Config Title:
  Example Backups Configuration
  #####################################
  Backup: App1
  Origin: /examples/dir1 | Destiny: /backups/dir1 | Retention: 45
  Backup: App2
  Origin: /examples/dir2 | Destiny: /backups/dir2 | Retention: 30
  Backup: App3
  Origin: /examples/dir3 | Destiny: /backups/dir3 | Retention: 30
  #####################################
  Checking directories:
  PASS: /examples/dir1
  PASS: /backups/dir1
  ++++ PASS!!!!!!! ++++
  #####################################
  Checking directories:
  PASS: /examples/dir2
  PASS: /backups/dir2
  ++++ PASS!!!!!!! ++++
  #####################################
  Checking directories:
  PASS: /examples/dir3
  PASS: /backups/dir3
  ++++ PASS!!!!!!! ++++
  #####################################
  Checking Retention for files
  #####################################
  /examples/dir1
  [file1 file2 file3]
  ====================================================
  /examples/dir2
  []
  nothing to backup in: /examples/dir2
  ====================================================
  /examples/dir3
  [file2 file3]
  ====================================================
  BACKING [{/examples/dir1 [file1 file2 file3] /backups/dir1} {/examples/dir3 [file2 file3] /backups/dir3}]
  Backup Successful
  Removed Original Files for /examples/dir1: [file1 file2 file3]
  Removed Original Files for /examples/dir3: [file2 file3]
