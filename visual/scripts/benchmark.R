data_all <- read.csv("visual/data/data_all_measurements.csv")

# scale from ns to ms
data_all$compress_duration_algorithm <-
  data_all$compress_duration_algorithm / 1000000

data_all$decompress_duration_algorithm <-
  data_all$decompress_duration_algorithm / 1000000

#head(data_all)
#data_all$algorithm

if (!require(ggplot2)) {install.packages("ggplot2")}

library(ggplot2)

algorithms <- unique(data_all$algorithm)

# box plot
matrics <- list(
                c(
                  "compress_duration_algorithm",
                  "compress_duration",
                  "Compression Duration [ms]"),
                c(
                  "decompress_duration_algorithm",
                  "decompress_duration",
                  "Decompression Duration [ms]"),
                c(
                  "compression_ratio",
                  "compression_ratio",
                  "Compression Ratio")) # Higher equals small resulting file

for (matric in matrics) {
  p <- ggplot(data_all, aes(x = algorithm, y = !!sym(matric[1]))) +
    geom_boxplot() +
    labs(
         title = paste(
                       "Min: ", min(data_all[[matric[1]]]),
                       "Max: ", max(data_all[[matric[1]]])),
         x = "Algorithm",
         y = matric[3])

  ggsave(
         filename = paste("boxplot_all_", matric[2], ".png"),
         plot = p,
         path = "visual/out/",
         create.dir = TRUE)

  # without BitFlipper, PPMd_exe, PPMonstr, and BZip2
  data <- data_all[!(data_all$algorithm %in% c("BitFlipper", "PPMd_exe", "PPMonstr", "BZip2")), ]
  p <- ggplot(data, aes(x = algorithm, y = !!sym(matric[1]))) +
    geom_boxplot() +
    labs(
         title = paste(
                       "Min: ", min(data[[matric[1]]]),
                       "Max: ", max(data[[matric[1]]])),
         x = "Algorithm",
         y = matric[3])

  ggsave(
         filename = paste("boxplot_all_without_slower_", matric[2], ".png"),
         plot = p,
         path = "visual/out/",
         create.dir = TRUE)

  for (algorithm in algorithms) {
    data <- data_all[data_all$algorithm == algorithm, ]

    min <- min(data[[matric[1]]])
    max <- max(data[[matric[1]]])
    mean <- mean(data[[matric[1]]])
    median <- median(data[[matric[1]]])

    p <- ggplot(data, aes(x = algorithm, y = !!sym(matric[1]))) +
      geom_boxplot() +
      labs(
           title = paste(
                         "Min: ", min,
                         "Max: ", max,
                         "Mean: ", mean,
                         "Median: ", median),
           x = "Algorithm",
           y = matric[3])

    ggsave(
           filename = paste("boxplot_", algorithm, "_", matric[2], ".png"),
           plot = p,
           path = "visual/out/")
  }
}

# scatter plot
data_scatter <- data.frame(
                           algorithm = character(),
                           compress_duration_algorithm = numeric(),
                           decompress_duration_algorithm = numeric(),
                           compression_ratio = numeric())

for (algorithm in algorithms) {
  data <- data_all[data_all$algorithm == algorithm, ]

  compress_duration <- median(data$compress_duration_algorithm)
  decompress_duration <- median(data$decompress_duration_algorithm)
  compression_ratio <- median(data$compression_ratio)

  data_scatter <- rbind(
    data_scatter,
    data.frame(
      algorithm = algorithm,
      compress_duration_algorithm = compress_duration,
      decompress_duration_algorithm = decompress_duration,
      compression_ratio = compression_ratio))
}

p <- ggplot(
            data_scatter,
            aes(
                x = compress_duration_algorithm,
                y = decompress_duration_algorithm,
                size = compression_ratio)) +
  geom_point(aes(color = algorithm), alpha = 0.7) +
  scale_size_continuous(range = c(10, 25)) +
  labs(
       title = "Compression vs Decompression Duration",
       x = "Compression Duration [ms]",
       y = "Decompression Duration [ms]")

ggsave(
       filename = "scatter_compress_decompress.png",
       plot = p,
       path = "visual/out/")

# multiple scatter plot compression and decompression speed vs size

data_all_without_bitFlipper <- data_all[data_all$algorithm != "BitFlipper", ]

data_all_without_bitFlipper$original_size_KB <- data_all_without_bitFlipper$original_size / 1000
data_all_without_bitFlipper$compress_speed_MB_s <- data_all_without_bitFlipper$original_size_KB / 1000 / ( data_all_without_bitFlipper$compress_duration_algorithm / 1000 )

p <- ggplot(
            data_all_without_bitFlipper,
            aes(
                x = original_size_KB,
                y = compress_speed_MB_s)) +
  geom_point(data = transform(data_all_without_bitFlipper, algorithm = NULL), colour = "grey85") +
  geom_point() +
  geom_smooth(method = "lm", colour = "black") +
  labs(
       title = "Compression Speed vs Original Size",
       x = "Original Size [KB]",
       y = "Compression Speed [MB/s]") +
  facet_wrap(vars(algorithm), nrow = 2)

ggsave(
       filename = "scatter_speeds_size.png",
       plot = p,
       path = "visual/out/")