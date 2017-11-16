# smallgopher
Smallworld GIS data web viewer.

- [GOSWORLD](#gosworld)
    - [About](#about)
    - [Run Cambridge example](#run-cambridge-example)
    - [Architecture](#architecture)
    - [Custom layers](#custom-layers)
    - [Custom styles](#custom-styles)

## About


This sample Web application to views data from Smallworld GIS in a browser.

![Cambridge example](https://github.com/kpawlik/smallgopher/blob/master/doc/cambridge2.png)

* Light web viewer for Smallworld GIS data
* 
* No installation/only executable no dependencies
* Small size
* Example for PNI and Cambridge db
* Customization for styles


## Run Cambridge example

1. Download repository from: https://github.com/kpawlik/smallgopher 
to: `C:\smallgopher`

2. Start HTTP server

```
"C:\smallgopher\cmd\smallgopher-server\runServer.cmd"
```

3. Load magik file to Smallworld session

```
load_file("C:\smallgopher\magik\goworld.magik")
```

4. run SW worker 

```
>>start_goworld_worker("w1", "C:\smallgopher\cmd\smallgopher-worker\smallgopher-worker.exe", 
":4001", "c:\tmp\log-w1.log")
$
```
5. Open Web browser nad go to address: `http://localhost:4000/`

![Cambridge example](https://github.com/kpawlik/smallgopher/blob/master/doc/cambridge1.png)

## Architecture

![Architecture example](https://github.com/kpawlik/smallgopher/blob/master/doc/smallgopher01.png)
## Custom layers

## Custom styles
