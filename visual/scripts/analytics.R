data_all <- read.csv("visual/data/data_all_analytics.csv")

if (!require(ggplot2)) {install.packages("ggplot2")}

library(ggplot2)

#head(data_all)


#min_value <- min(data_all$compressed_size_ratio)

#min_value

min_row <- data_all[which.min(data_all$compressed_size_ratio), ]

min_row


p <- ggplot(data_all[data_all$file_name == "HE2GQLBH2WBUUJKRBNZCYY5QZWTVCY35.pdf", ], aes(x = window_length_bytes, y = compressed_size_ratio)) +
  geom_line() +
  geom_point() +
  labs(title = "Compressed Size Ratio vs Window Length Bytes",
       x = "Window Length [Bytes]",
       y = "Compressed Size Ratio") +
  ylim(0, 100)

ggsave(
         filename = paste("line_compressed_size_ratio_", "HE2GQLBH2WBUUJKRBNZCYY5QZWTVCY35.pdf", ".png"),
         plot = p,
         path = "visual/out/",
         create.dir = TRUE)


p <- ggplot(data_all[data_all$file_name == "HE2GQLBH2WBUUJKRBNZCYY5QZWTVCY35.pdf", ], aes(x = window_length_bytes, y = compressed_size_ratio)) +
  geom_point() +
  labs(title = "Compressed Size Ratio vs Window Length Bytes",
       x = "Window Length [Bytes]",
       y = "Compressed Size Ratio") +
  ylim(0, 100)

ggsave(
         filename = paste("scatter_compressed_size_ratio_", "HE2GQLBH2WBUUJKRBNZCYY5QZWTVCY35.pdf", ".png"),
         plot = p,
         path = "visual/out/",
         create.dir = TRUE)


file_names <- unique(data_all$file_name)

min_values <- data.frame(file_name = character(), window_length_bytes = numeric(), compressed_size_ratio = numeric())

for (file_name in file_names) {
  file_data <- data_all[data_all$file_name == file_name, ]

  min_value <- file_data[which.min(file_data$compressed_size_ratio), ]

  if (min_value$compressed_size_ratio == 100) {
    next
  }

  min_values <- rbind(min_values, min_value)
}


p <- ggplot(min_values, aes(x = window_length_bytes, y = compressed_size_ratio)) +
  geom_point() +
  labs(title = paste(
        "Min:", min(min_values$compressed_size_ratio),
        "Max:", max(min_values$compressed_size_ratio),
        "File count:", nrow(min_values), "/", length(file_names)),
       x = "Window Length [Bytes]",
       y = "Compressed Size Ratio") +
  ylim(0, 100)

ggsave(
        filename = paste("scatter_min_compressed_size_ratio", ".png"),
        plot = p,
        path = "visual/out/",
        create.dir = TRUE)