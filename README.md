# Stopwatch

## Import Teams

```
$ rm tmp/2023/teams/XS.yaml && make install && stopwatch2 import-teams tmp/2023/teams/XS.csv tmp/2023/teams/XS.yaml



```

## Record Race

```
$ stopwatch2 shell tmp/2023/races/XS.yaml
**********************************
type 'start' when the race begins!
**********************************
⏱ start
⏱ 1
⏱ 19
⏱ 2
⏱ exit
bye! 👋
```

## Generate Results

```
$ stopwatch2 generate-report "XS 2023" tmp/2023/teams/XS.yaml tmp/2023/races/XS.yaml tmp/2023/reports

$ asciidoctor-pdf tmp/2023/reports/*.adoc 
```