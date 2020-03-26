Predicted Target Profile towards proteins of interest for Covid-19
==================================================================

This is a project aimed at running a similar pipeline like in the [PTP
project](https://github.com/pharmbio/ptp-project) on data for the Covid-19
disease, caused by the Coronavirus.

The basic idea is to use the 67 druggable human protein targets reported in
[this study](https://doi.org/10.1101/2020.03.22.002386) to develop predictive
models for them, based on binding data for those proteins in open datasets.


## Requirements

To run the scripts and workflows in this repository, you need the following
software installed.

- A unix-like environment (Ubuntu, MacOS, or [Windows with MSYS2](https://rillabs.com/posts/linux-like-non-wsl-terminal-env-on-windows))
- The bash shell
- [Go 1.5 or later](https://golang.org/)
- [SciPipe 0.9.6 or later](https://scipipe.org/)
- The `ssconvert` command from [Gnumeric](http://www.gnumeric.org/)

## How to reproduce

#### Easy way

To download all data and reproduce all experiments, run the `runall.sh` shell script:

```bash
./runall.sh
```

... or:

```bash
bash runall.sh
```

#### Manually with more control

If you want to run the individual steps manually to have more control over
what's happening, then follow the following steps.

- Get raw data:

  ```bash
  go run getrawdata.go
  ```

- Execute the counting experiment:

  ```bash
  cd exp/20200525-count
  go run count.go
  ```
