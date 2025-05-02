data_all <- read.csv("visual/data/data_all_cloud.csv")

if (!require(ggplot2)) {install.packages("ggplot2")}

library(ggplot2)

# convert ns to ms
data_all$duration_upload <- data_all$duration_upload / 1000000
data_all$duration_download <- data_all$duration_download / 1000000

# boxplot tier_name vs duration_upload
p <- ggplot(data_all, aes(x = tier_name, y = duration_upload)) +
  geom_boxplot() +
  labs(title = "Boxplot of Duration Upload by Tier Name",
       x = "Tier Name",
       y = "Duration Upload [ms]")

ggsave(
       filename = "boxplot_tier_name_duration_upload.png",
       plot = p + ylim(0, max(data_all$duration_upload)),
       path = "visual/out/")

ggsave(
       filename = "boxplot_tier_name_duration_upload_log10.png",
       plot = p + scale_y_log10(),
       path = "visual/out/")

# boxplot tier_name vs duration_upload
p <- ggplot(data_all, aes(x = tier_name, y = duration_download)) +
  geom_boxplot() +
  labs(title = "Boxplot of Duration Download by Tier Name",
       x = "Tier Name",
       y = "Duration Download [ms]")

ggsave(
       filename = "boxplot_tier_name_duration_download.png",
       plot = p + ylim(0, max(data_all$duration_download)),
       path = "visual/out/")

ggsave(
       filename = "boxplot_tier_name_duration_download_log10.png",
       plot = p + scale_y_log10(),
       path = "visual/out/")

# scatter plot duration_upload vs size
p <- ggplot(data_all, aes(x = size, y = duration_upload)) +
  labs(title = "Scatter Plot of Duration Upload vs Size",
       x = "Size [bytes]",
       y = "Duration Upload [ms]")

ggsave(
        filename = "scatter_duration_upload_vs_size.png",
        plot = p + geom_point() + ylim(0, max(data_all$duration_upload)),
        path = "visual/out/")

ggsave(
        filename = "scatter_duration_upload_vs_size_line.png",
        plot = p + geom_point(alpha = 0.5, colour = "grey") + ylim(0, max(data_all$duration_upload)) + geom_smooth(method = "lm", colour = "black"),
        path = "visual/out/")

# scatter plot duration_download vs size
p <- ggplot(data_all, aes(x = size, y = duration_download)) +
  labs(title = "Scatter Plot of Duration Download vs Size",
       x = "Size [bytes]",
       y = "Duration Download [ms]")

ggsave(
        filename = "scatter_duration_download_vs_size.png",
        plot = p + geom_point() + ylim(0, max(data_all$duration_download)),
        path = "visual/out/")

ggsave(
        filename = "scatter_duration_download_vs_size_line.png",
        plot = p + geom_point(alpha = 0.5, colour = "grey") + ylim(0, max(data_all$duration_download)) + geom_smooth(method = "lm", colour = "black"),
        path = "visual/out/")

# scatter plot duration_upload vs duration_download
p <- ggplot(data_all, aes(x = duration_upload, y = duration_download)) +
  geom_point() +
  labs(title = "Scatter Plot of Duration Upload vs Duration Download",
       x = "Duration Upload [ms]",
       y = "Duration Download [ms]")

ggsave(
        filename = "scatter_duration_upload_vs_duration_download.png",
        plot = p,
        path = "visual/out/")