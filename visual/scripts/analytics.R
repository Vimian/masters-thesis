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
       y = "Compressed Size Ratio") +
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
  y = "Compressed Size Ratio")

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
  y = "Compressed Size Ratio")

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

# scatter plot of file size vs window length bytes
#p <- ggplot(min_values, aes(x = window_length_bytes, y = file_size)) +
#  geom_point() +
#  geom_smooth(method = "lm", colour = "black") +
#  labs(title = paste(
#        "Min:", min(min_values$file_size),
#        "Max:", max(min_values$file_size),
#        "File count:", nrow(min_values), "/", length(file_names)),
#       x = "Window Length [Bytes]",
#       y = "File Size [Bytes]") +
#  ylim(0, max(min_values$file_size))
#
#ggsave(
#        filename = paste("scatter_file_size", ".png"),
#        plot = p,
#        path = "visual/out/",
#        create.dir = TRUE)