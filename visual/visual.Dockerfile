FROM r-base:4.4.2

RUN apt-get update && apt-get install -y libpq-dev libssl-dev

RUN Rscript -e "install.packages(c('DBI', 'RPostgres', 'ggplot2'))"