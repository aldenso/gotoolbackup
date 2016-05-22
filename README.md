gotoolbackup
============

Program to create backups using toml file, where you indicate origin/destiny directories and retention period in days.

Note: this is not to replace shell scripts or some enterprise tool.

Usage

    > ./gotoolbackup -help
    Usage of ./gotoolbackup:
      -log string
            indicate the log name pattern. (default "gotoolbackup")
      -remove
            indicate if you want to remove original files after backup.
      -tomlfile string
            indicate tomlfile to read backups details. (default "gobackup.toml")


    > ./gotoolbackup -log=output -tomlfile=/tmp/gobackup.toml -remove
    #### Running with values ####
    tomlfile: /tmp/gobackup.toml
    remove: true
    log: output
    #####################################
    Config Title:
    Example Backups Configuration
    #####################################
    Backup: App2
    Origin: /examples/dir2 | Destiny: /backups/dir2 | Retention: 30
    Backup: App3
    Origin: /examples/dir3 | Destiny: /backups/dir3 | Retention: 30
    Backup: App1
    Origin: /examples/dir1 | Destiny: /backups/dir1 | Retention: 45
    #####################################
    Checking directories:
    PASS: /examples/dir3
    PASS: /backups/dir3
    ++++ PASS!!!!!!! ++++
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
    Running backups for:
    /examples/dir1: file1,file2,file3 size in bytes: 159223635
    /examples/dir3: file2,file3 size in bytes: 47185920
    Backup Successful
    Removed Original Files for /examples/dir1: [file1 file2 file3]
    Removed Original Files for /examples/dir3: [file2 file3]
    old files removed
    gotoolbackup ended! in: 3.258203627s
