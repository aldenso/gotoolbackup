gotoolbackup
============

Program to create backups using toml file, where you indicate origin/destiny directories and retention period in days.

Note: this is not created to replace shell scripts or some enterprise tool.

Usage

```
$ ./gotoolbackup --help
Usage of ./gotoolbackup:
  -log string
        indicate the log name pattern. (default "gotoolbackup")
  -remove
        indicate if you want to remove original files after backup.
  -tomlfile string
        indicate tomlfile to read backups details. (default "gobackup.toml")
```

Before the backup this is the content.

```
$ ls -lrthR /examples/
/examples/:
total 12K
drwxr-xr-x. 2 aldenso aldenso 4.0K Sep 24 20:19 dir1
drwxr-xr-x. 2 aldenso aldenso 4.0K Sep 24 20:19 dir2
drwxr-xr-x. 2 aldenso aldenso 4.0K Sep 24 20:19 dir3

/examples/dir1:
total 65M
-rw-r--r--. 1 aldenso aldenso  41M Sep 10  2014 file1.txt
-rw-r--r--. 1 aldenso aldenso  10M Oct 12  2015 file3.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Oct 12  2015 file2.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:48 file1_clone.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Sep 24 14:48 file2_clone.txt
-rw-r--r--. 1 aldenso aldenso  10M Sep 24 14:48 file3_clone.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:49 file1_root.txt
-rw-rw-r--. 1 aldenso aldenso 3.8M Sep 24 15:00 original.tar.gz

/examples/dir2:
total 66M
-rw-r--r--. 1 aldenso aldenso  41M Sep 10  2015 file4.txt
-rw-r--r--. 1 aldenso aldenso  10M Oct 12  2015 file3.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Oct 12  2015 file2.txt
-rw-r--r--. 1 aldenso aldenso 330K Oct 12  2015 file1.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:48 file1_clone.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Sep 24 14:48 file2_clone.txt
-rw-r--r--. 1 aldenso aldenso  10M Sep 24 14:48 file3_clone.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:49 file1_root.txt
-rw-r--r--. 1 root root 3.8M Sep 24 15:05 original.tar.gz

/examples/dir3:
total 66M
-rw-r--r--. 1 aldenso aldenso  41M Sep 10  2015 file4.txt
-rw-r--r--. 1 aldenso aldenso  10M Oct 12  2015 file3.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Oct 12  2015 file2.txt
-rw-r--r--. 1 aldenso aldenso 330K Oct 12  2015 file1.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:48 file1_clone.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Sep 24 14:48 file2_clone.txt
-rw-r--r--. 1 aldenso aldenso  10M Sep 24 14:48 file3_clone.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:49 file1_root.txt
-rw-r--r--. 1 root root 3.8M Sep 24 15:05 original.tar.gz
```

Now we run the backup with specific options.

```
$ ./gotoolbackup -log=output -tomlfile=/tmp/gobackup.toml -remove
#### Running with values ####
tomlfile: /tmp/gobackup.toml
remove: true
log: output
Reading tomlfile: /tmp/gobackup.toml
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
[file1.txt file2.txt file3.txt]
====================================================
/examples/dir2
[file1.txt file2.txt file3.txt file4.txt]
====================================================
/examples/dir3
[file1.txt file2.txt file3.txt file4.txt]
====================================================
Running backups for:
/examples/dir1: file1.txt,file2.txt,file3.txt - size in bytes: 52918656
/examples/dir2: file1.txt,file2.txt,file3.txt,file4.txt - size in bytes: 53255983
/examples/dir3: file1.txt,file2.txt,file3.txt,file4.txt - size in bytes: 53255983
backup file: /backups/dir2/backup_2016-09-24T201523-0400.tar.gz - size in bytes: 3819530
backup file: /backups/dir3/backup_2016-09-24T201523-0400.tar.gz - size in bytes: 3819530
backup file: /backups/dir1/backup_2016-09-24T201523-0400.tar.gz - size in bytes: 3809546
Backup Successful
Removed Original Files for /examples/dir1: [file1.txt file2.txt file3.txt]
Removed Original Files for /examples/dir2: [file1.txt file2.txt file3.txt file4.txt]
Removed Original Files for /examples/dir3: [file1.txt file2.txt file3.txt file4.txt]
old files removed
gotoolbackup ended! in: 14.468081137s
```

Let's check the log.

```
$ cat output_2016-09-24T201523-0400.log
gotoolbackup: 2016/09/24 20:15:23 Reading tomlfile: /tmp/gobackup.toml
gotoolbackup: 2016/09/24 20:15:23 Checking Retention for files
gotoolbackup: 2016/09/24 20:15:23 Running backups for:
gotoolbackup: 2016/09/24 20:15:23 /examples/dir1: file1.txt,file2.txt,file3.txt - size in bytes: 52918656
gotoolbackup: 2016/09/24 20:15:23 /examples/dir2: file1.txt,file2.txt,file3.txt,file4.txt - size in bytes: 53255983
gotoolbackup: 2016/09/24 20:15:23 /examples/dir3: file1.txt,file2.txt,file3.txt,file4.txt - size in bytes: 53255983
gotoolbackup: 2016/09/24 20:15:37 backup file: /backups/dir2/backup_2016-09-24T201523-0400.tar.gz - size in bytes: 3819530
gotoolbackup: 2016/09/24 20:15:37 backup file: /backups/dir3/backup_2016-09-24T201523-0400.tar.gz - size in bytes: 3819530
gotoolbackup: 2016/09/24 20:15:37 backup file: /backups/dir1/backup_2016-09-24T201523-0400.tar.gz - size in bytes: 3809546
gotoolbackup: 2016/09/24 20:15:37 Backup Successful
gotoolbackup: 2016/09/24 20:15:37 old files removed
gotoolbackup: 2016/09/24 20:15:37 gotoolbackup ended! in: 14.468081137s
```

This is the content after the backup with remove option of the original files.

```
$ ls -lrthR /examples/
/examples/:
total 12K
drwxr-xr-x. 2 aldenso aldenso 4.0K Sep 24 20:23 dir1
drwxr-xr-x. 2 aldenso aldenso 4.0K Sep 24 20:23 dir2
drwxr-xr-x. 2 aldenso aldenso 4.0K Sep 24 20:23 dir3

/examples/dir1:
total 15M
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:48 file1_clone.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Sep 24 14:48 file2_clone.txt
-rw-r--r--. 1 aldenso aldenso  10M Sep 24 14:48 file3_clone.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:49 file1_root.txt
-rw-rw-r--. 1 aldenso aldenso 3.8M Sep 24 15:00 original.tar.gz

/examples/dir2:
total 15M
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:48 file1_clone.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Sep 24 14:48 file2_clone.txt
-rw-r--r--. 1 aldenso aldenso  10M Sep 24 14:48 file3_clone.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:49 file1_root.txt
-rw-r--r--. 1 root root 3.8M Sep 24 15:05 original.tar.gz

/examples/dir3:
total 15M
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:48 file1_clone.txt
-rw-r--r--. 1 aldenso aldenso 4.4K Sep 24 14:48 file2_clone.txt
-rw-r--r--. 1 aldenso aldenso  10M Sep 24 14:48 file3_clone.txt
-rw-r--r--. 1 aldenso aldenso 330K Sep 24 14:49 file1_root.txt
-rw-r--r--. 1 root root 3.8M Sep 24 15:05 original.tar.gz
```

Finally let's check the backups.

```
$ ls -lrthR /backups/
/backups/:
total 0
drwxr-xr-x. 2 aldenso aldenso 49 Sep 24 20:25 dir1
drwxr-xr-x. 2 aldenso aldenso 49 Sep 24 20:26 dir3
drwxr-xr-x. 2 aldenso aldenso 49 Sep 24 20:26 dir2

/backups/dir1:
total 3.7M
-rw-rw-r--. 1 aldenso aldenso 3.7M Sep 24 20:15 backup_2016-09-24T201523-0400.tar.gz

/backups/dir3:
total 3.7M
-rw-rw-r--. 1 aldenso aldenso 3.7M Sep 24 20:15 backup_2016-09-24T201523-0400.tar.gz

/backups/dir2:
total 3.7M
-rw-rw-r--. 1 aldenso aldenso 3.7M Sep 24 20:15 backup_2016-09-24T201523-0400.tar.gz

$ tar tvf /backups/dir1/backup_2016-09-24T201523-0400.tar.gz
-rw-r--r-- 1000/1000  42428422 2014-09-10 12:30 file1.txt
-rw-r--r-- 1000/1000      4474 2015-10-12 07:30 file2.txt
-rw-r--r-- 1000/1000  10485760 2015-10-12 07:30 file3.txt

$ tar tvf /backups/dir2/backup_2016-09-24T201523-0400.tar.gz
-rw-r--r-- 1000/1000    337327 2015-10-12 07:30 file1.txt
-rw-r--r-- 1000/1000      4474 2015-10-12 07:30 file2.txt
-rw-r--r-- 1000/1000  10485760 2015-10-12 07:30 file3.txt
-rw-r--r-- 1000/1000  42428422 2015-09-10 09:15 file4.txt

$ tar tvf /backups/dir3/backup_2016-09-24T201523-0400.tar.gz
-rw-r--r-- 1000/1000    337327 2015-10-12 07:30 file1.txt
-rw-r--r-- 1000/1000      4474 2015-10-12 07:30 file2.txt
-rw-r--r-- 1000/1000  10485760 2015-10-12 07:30 file3.txt
-rw-r--r-- 1000/1000  42428422 2015-09-10 09:15 file4.txt
```
