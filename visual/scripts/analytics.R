data_all <- read.csv("visual/data/data_all_analytics.csv")

if (!require(ggplot2)) {install.packages("ggplot2")}

library(ggplot2)

#head(data_all)


#min_value <- min(data_all$compressed_size_ratio)

#min_value

min_row <- data_all[which.min(data_all$compressed_size_ratio), ]

min_row

p <- ggplot(
            data_all[data_all$file_name ==
                       "HE2GQLBH2WBUUJKRBNZCYY5QZWTVCY35.pdf" &
                       data_all$compressed_size_ratio != 100, ],
            aes(x = window_length_bytes, y = compressed_size_ratio)) +
  geom_point() +
  labs(title = "Compressed Size Ratio vs Window Length Bytes",
       x = "Window Length [Bytes]",
       y = "Compressed Size Ratio [%]") +
  ylim(min(min_row$compressed_size_ratio), 100)

ggsave(
       filename = paste("HE2GQLBH2WBUUJKRBNZCYY5QZWTVCY35.pdf", ".png"),
       plot = p,
       path = "visual/out/",
       create.dir = TRUE)


file_names <- unique(data_all$file_name)

min_values <- data.frame()

for (file_name in file_names) {
  file_data <- data_all[data_all$file_name == file_name, ]

  min_value <- file_data[which.min(file_data$compressed_size_ratio), ]

  if (min_value$compressed_size_ratio == 100) {
    next
  }

  min_values <- rbind(min_values, min_value)
}

# scatter plot of compressed size ratio vs window length bytes
p <- ggplot(min_values, aes(
                            x = window_length_bytes,
                            y = compressed_size_ratio)) +
  geom_point() +
  geom_smooth(method = "lm", colour = "black") +
  labs(title = paste(
                     "Min:", min(min_values$compressed_size_ratio),
                     "Max:", max(min_values$compressed_size_ratio),
                     "File count:", nrow(min_values), "/", length(file_names)),
  x = "Window Length [Bytes]",
  y = "Compressed Size Ratio [%]")

ggsave( # 0 to 100
       filename = paste("scatter_min_compressed_size_ratio_0_100", ".png"),
       plot = p + ylim(0, 100),
       path = "visual/out/",
       create.dir = TRUE)

ggsave( # scaled
       filename = paste("scatter_min_compressed_size_ratio_scaled", ".png"),
       plot = p + ylim(min(min_values$compressed_size_ratio), 100),
       path = "visual/out/",
       create.dir = TRUE)

# scatter plot of compressed size ratio vs file size
p <- ggplot(min_values, aes(x = file_size, y = compressed_size_ratio)) +
  geom_point() +
  labs(title = paste(
                     "Min:", min(min_values$compressed_size_ratio),
                     "Max:", max(min_values$compressed_size_ratio),
                     "File count:", nrow(min_values), "/", length(file_names)),
  x = "File Size [Bytes]",
  y = "Compressed Size Ratio [%]")

ggsave( # 0 to 100
       filename = paste("scatter_compressed_size_ratio_0_100", ".png"),
       plot = p + geom_smooth(method = "lm", colour = "black") +
         ylim(0, 100),
       path = "visual/out/",
       create.dir = TRUE)

ggsave( # scaled
       filename = paste("scatter_compressed_size_ratio_scaled", ".png"),
       plot = p + geom_smooth(method = "lm", colour = "black") +
         ylim(min(min_values$compressed_size_ratio), 100),
       path = "visual/out/",
       create.dir = TRUE)

ggsave( # scaled with exponential trendline
       filename = paste("scatter_compressed_size_ratio_scaled_exponential",
                        ".png"),
       plot = p + geom_smooth(method = "lm", colour = "black",
                              formula = y ~ poly(x, 2)) +
         ylim(min(min_values$compressed_size_ratio), 100),
       path = "visual/out/",
       create.dir = TRUE)

# scatter plot of window length bytes vs file size
p <- ggplot(min_values, aes(x = file_size, y = window_length_bytes)) +
  geom_point() +
  geom_smooth(method = "lm", colour = "black") +
  labs(title = paste(
                     "Min:", min(min_values$window_length_bytes),
                     "Max:", max(min_values$window_length_bytes),
                     "File count:", nrow(min_values), "/", length(file_names)),
  x = "File Size [Bytes]",
  y = "Window Length [Bytes]")

ggsave( # 0 to max
       filename = paste("scatter_window_length_bytes_0_max", ".png"),
       plot = p + ylim(0, max(data_all$window_length_bytes)),
       path = "visual/out/",
       create.dir = TRUE)

ggsave( # scaled
       filename = paste("scatter_window_length_bytes_scaled", ".png"),
       plot = p + ylim(0, max(min_values$window_length_bytes)),
       path = "visual/out/",
       create.dir = TRUE)






window_length_bytes_count <- data.frame(
  window_length_bytes = seq(
                            min(data_all$window_length_bytes),
                            max(data_all$window_length_bytes), 1),
  count = rep(0, length(seq(
                            min(data_all$window_length_bytes),
                            max(data_all$window_length_bytes), 1)))
)

for (i in 1:nrow(data_all)) {
  data_point <- data_all[i, ]

  if (data_point$compressed_size_ratio >= 100) {
    next
  }

  window_length_bytes_count$count[data_point$window_length_bytes] <-
    window_length_bytes_count$count[data_point$window_length_bytes] + 1
}

times_all_files_success_compressed <-
  window_length_bytes_count[
                            window_length_bytes_count$count ==
                            length(file_names), ]

# scatter plot of window_length_bytes vs count
p <- ggplot(window_length_bytes_count[window_length_bytes_count$count != 0, ],
            aes(x = window_length_bytes, y = count)) +
  geom_point() +
  labs(title = paste(
                     "Min:", min(window_length_bytes_count$count),
                     "Max:", max(window_length_bytes_count$count),
                     "Times all files got more compact:",
                     nrow(times_all_files_success_compressed),
                     "/", max(data_all$window_length_bytes)),
  x = "Window Length [Bytes]",
  y = "Count")

ggsave(
       filename = paste("scatter_plot_window_length_bytes_count", ".png"),
       plot = p,
       path = "visual/out/",
       create.dir = TRUE)

ggsave(
       filename = paste("scatter_plot_window_length_bytes_count_0_1500", ".png"),
       plot = p + xlim(0, 1500),
       path = "visual/out/",
       create.dir = TRUE)

ggsave(
       filename = paste("scatter_plot_window_length_bytes_count_0_110", ".png"),
       plot = p + xlim(0, 110),
       path = "visual/out/",
       create.dir = TRUE)

window_length_bytes_count_sorted <-
  window_length_bytes_count[
                            order(window_length_bytes_count$count,
                                  decreasing = TRUE), ]

write.csv(
          window_length_bytes_count_sorted[seq(1, 100), ],
          file = "visual/out/window_length_bytes_count_sorted_0_100.csv",
          row.names = TRUE)

write.csv(
          window_length_bytes_count[seq(101, 1010, 101), ],
          file = "visual/out/window_length_bytes_count_101_1010.csv",
          row.names = TRUE)

data_best_window_length_bytes <-
  data_all[
    data_all$window_length_bytes == 40 |
    data_all$window_length_bytes == 44 |
    data_all$window_length_bytes == 48 |
    data_all$window_length_bytes == 42 |
    data_all$window_length_bytes == 38 |
    data_all$window_length_bytes == 46 |
    data_all$window_length_bytes == 52,
  ]

# boxplot of compressed_size_ratio of the best window_length_bytes
p <- ggplot(
            data_best_window_length_bytes,
            aes(x = factor(window_length_bytes), y = compressed_size_ratio)) +
  geom_boxplot() +
  labs(title = paste(
                     "Min:",
                     min(data_best_window_length_bytes$compressed_size_ratio),
                     "Max:",
                     max(data_best_window_length_bytes$compressed_size_ratio)),
  x = "Window Length [Bytes]",
  y = "Compressed Size Ratio [%]")

ggsave(
       filename = paste("boxplot_compressed_size_ratio", ".png"),
       plot = p,
       path = "visual/out/",
       create.dir = TRUE)

data_best_window_length_bytes_only <-
  data_best_window_length_bytes[
    data_best_window_length_bytes$compressed_size_ratio < 100,
  ]

# boxplot of compressed_size_ratio of the best window_length_bytes
p <- ggplot(
            data_best_window_length_bytes_only,
            aes(x = factor(window_length_bytes), y = compressed_size_ratio)) +
  geom_boxplot() +
  labs(title = paste(
                     "Min:",
                     min(data_best_window_length_bytes_only$compressed_size_ratio),
                     "Max:",
                     max(data_best_window_length_bytes_only$compressed_size_ratio)),
  x = "Window Length [Bytes]",
  y = "Compressed Size Ratio [%]")

ggsave(
       filename = paste("boxplot_compressed_size_ratio_only_compressed", ".png"),
       plot = p,
       path = "visual/out/",
       create.dir = TRUE)

data_what_window_length_bytes <-
  data_all[
    data_all$window_length_bytes == 101 |
    data_all$window_length_bytes == 202 |
    data_all$window_length_bytes == 303 |
    data_all$window_length_bytes == 404 |
    data_all$window_length_bytes == 505 |
    data_all$window_length_bytes == 606 |
    data_all$window_length_bytes == 707 |
    data_all$window_length_bytes == 808 |
    data_all$window_length_bytes == 909 |
    data_all$window_length_bytes == 1010,
  ]

# boxplot of compressed_size_ratio of the best window_length_bytes
p <- ggplot(
            data_what_window_length_bytes,
            aes(x = factor(window_length_bytes), y = compressed_size_ratio)) +
  geom_boxplot() +
  labs(title = paste(
                     "Min:",
                     min(data_what_window_length_bytes$compressed_size_ratio),
                     "Max:",
                     max(data_what_window_length_bytes$compressed_size_ratio)),
  x = "Window Length [Bytes]",
  y = "Compressed Size Ratio [%]")

ggsave(
       filename = paste("boxplot_compressed_size_ratio_what", ".png"),
       plot = p,
       path = "visual/out/",
       create.dir = TRUE)

data_what_window_length_bytes_only <-
  data_what_window_length_bytes[
    data_what_window_length_bytes$compressed_size_ratio < 100,
  ]

# boxplot of compressed_size_ratio of the best window_length_bytes
p <- ggplot(
            data_what_window_length_bytes_only,
            aes(x = factor(window_length_bytes), y = compressed_size_ratio)) +
  geom_boxplot() +
  labs(title = paste(
                     "Min:",
                     min(data_what_window_length_bytes_only$compressed_size_ratio),
                     "Max:",
                     max(data_what_window_length_bytes_only$compressed_size_ratio)),
  x = "Window Length [Bytes]",
  y = "Compressed Size Ratio [%]")

ggsave(
       filename = paste("boxplot_compressed_size_ratio_what_only_compressed", ".png"),
       plot = p,
       path = "visual/out/",
       create.dir = TRUE)
