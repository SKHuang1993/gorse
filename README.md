<img width=160 src="https://img.sine-x.com/gorse.png"/>

# gorse: Go Recommender System Engine

| Build | Build (AVX2) | Coverage | Document | Report |
|---|---|---|---|---|
| [![Build Status](https://travis-matrix-badges.herokuapp.com/repos/zhenghaoz/gorse/branches/master/1)](https://travis-ci.org/zhenghaoz/gorse) | [![Build Status](https://travis-matrix-badges.herokuapp.com/repos/zhenghaoz/gorse/branches/master/2)](https://travis-ci.org/zhenghaoz/gorse) | [![codecov](https://codecov.io/gh/zhenghaoz/gorse/branch/master/graph/badge.svg)](https://codecov.io/gh/zhenghaoz/gorse) | [![GoDoc](https://godoc.org/github.com/zhenghaoz/gorse?status.svg)](https://godoc.org/github.com/zhenghaoz/gorse) | [![Go Report Card](https://goreportcard.com/badge/github.com/zhenghaoz/gorse)](https://goreportcard.com/report/github.com/zhenghaoz/gorse) |

`gorse` is a a transparent recommender system engine over SQL database based on collaborative filtering written in Go.

## Install

```bash
go get github.com/zhenghaoz/gorse/...
```

## Build

```bash
go build github.com/zhenghaoz/gorse/cmd/gorse
```

If the CPU of your device supports AVX2 and FMA3 instructions, use the `avx2` build tag to enable AVX2 support.

```bash
go build -tags='avx2' github.com/zhenghaoz/gorse/cmd/gorse
```

## Usage

It's easy to setup a recomendation service with `gorse`. 

- **Step 1**: initialize the database.

```bash
./gorse init user:pass@host/database
```

It connects to a SQL database and creates several tables.  `user:pass@host/database` is the database used to store data of `gorse`.

- **Step 2**: Import ratings and Items.

```
./gorse data user:pass@host/database \
	--import-ratings-csv u.data \
	--import-items-csv u.item
```

It imports ratings and items from CSV files. `u.data` is the CSV file of ratings in MovieLens 100K dataset and `u.item` is the CSV file of items in MovieLens 100K dataset.

- **Step 3**: Start a server.

```bash
./gorse server -c config.toml
```

Load configurations from `config.toml` and start a recommendation server. It may take a while to generate all recommendations.

```toml
# This section declares settings for the server.
[server]
host = "127.0.0.1"      # server host
port = 8080             # server port

# This section declares setting for the database.
[database]
driver = "mysql"        # database driver
access = "gorse:password@/gorse"# database access

# This section declares settings for recommendation.
[recommend]
model = "svd"           # recommendation model
cache_size = 100        # the number of cached recommendations
update_threshold = 10   # update model when more than 10 ratings are added
check_period = 1        # check for update every one minute
similarity = "pearson"  # similarity metric for neighbors

# This section declares hyperparameters for the recommendation model.
[params]
optimizer = "bpr"       # the optimizer to oprimize matrix factorization model
n_factors = 10          # the number of latent factors
reg = 0.01              # regularization strength
lr = 0.05               # learning rate
n_epochs = 100          # the number of learning epochs
init_mean = 0.0         # the mean of initial latent factors initilaized by Gaussian distribution
init_std = 0.001        # the standard deviation of initial latent factors initilaized by Gaussian distribution
```

- **Step 4**: Send requests.

```bash
curl 127.0.0.1:8080/recommends/1?number=5
```

Request 5 recommended items for the 1-th user. See [more APIs](https://github.com/zhenghaoz/gorse/wiki/RESTful-APIs) in wiki.

```json
{
  "Failed": false,
  "Items": [284, 448, 763, 276, 313]
}
```

## Document

- Visit [GoDoc](https://godoc.org/github.com/zhenghaoz/gorse) for detailed documentation of codes.
- Visit [Wiki](https://github.com/zhenghaoz/gorse/wiki) for tutorial, examples and high-level introduction.

## Performance

gorse is much faster than Surprise, and comparable to librec while using less memory space than Surprise and librec. The memory efficiency is achieved by sophisticated data structures. 

- cross-validation of SVD on MovieLens 100K:

<img width=320 src="https://img.sine-x.com/perf_time_svd_ml_100k.png"><img width=320 src="https://img.sine-x.com/perf_mem_svd_ml_100k.png">

- cross-validation of SVD on MovieLens 1M:

<img width=320 src="https://img.sine-x.com/perf_time_svd_ml_1m.png"><img width=320 src="https://img.sine-x.com/perf_mem_svd_ml_1m.png">

## Contributors

[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/0)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/0)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/1)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/1)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/2)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/2)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/3)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/3)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/4)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/4)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/5)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/5)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/6)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/6)[![](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/images/7)](https://sourcerer.io/fame/zhenghaoz/zhenghaoz/gorse/links/7)

Any kind of contribution is expected: report a bug, give a advice or even create a pull request.

## Acknowledgments

`gorse` was inspired by following projects:

- [Guibing Guo's librec](https://github.com/guoguibing/librec)
- [Nicolas Hug's Surprise](https://github.com/NicolasHug/Surprise)
- [Golang Samples's gopher-vector](https://github.com/golang-samples/gopher-vector)

## Limitations

`gorse` has limitations and might not be applicable to some scenarios:

- **No Scalability**: Since `gorse` is a recommendation service on a single host, it's unable to handle large data. The bottleneck might be memory size or the performance of SQL database.
- **No Feature Engineering**:  `gorse` only uses interactions between items and users while features of items, users and contexts are ignored. This is not the best practice in real world.
- **Naive Policy**: There are lots of considerations on a recommender system such as the freshness of items, the variation of users' preferences ,etc. They are not included in `gorse`. 
