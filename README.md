# Gophorganiser
This tool allows you to organise all pictures you have in different directories into a target directory, ordered by datetime.

## Why this software?
You have a phone, a tablet, a camera, and when you synchronise/backup them all (which I hope you do) to a hard disk (NAS, cloud, etc), every device has its own folder. You spend a lot of time moving pictures around, trying to make everything like a timeline: pictures of the same day from all your devices, please show up together.

## How it works?
First of all, this command line tool needs access to the disks where you store your pictures, and to the target disk where you want to have them organised by datetime.

To accomplish that, the CLI:

1. will ask for a source directory.
1. will ask for another directory, until the user answers `"No"`.
1. will ask for a target directory, where all pictures will be copied.
1. once a directory is processed (pictures removed and copied), it will remove source directory contents.
1. the target/moved file name will consist in: timestamp of the original picture plus `_` plus the original name and extension. I.e., the file `DCN-001.jpg` will be processed as `20191222_DCN-001.jpg`

>A directory will be considered processed when all its contents (photos and videos) are copied to the target directory.