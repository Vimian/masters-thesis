pricelist <- list(
        c(
                "S3 Standard",
                2163.71,
                0.3565,
                0.02852
        ),
        c(
                "S3 Intelligent FA",
                2163.71,
                0.3565,
                0.02852
        ),
        c(
                "S3 Intelligent IA",
                1175.928,
                0.3565,
                0.02852
        ),
        c(
                "S3 Intelligent AIA",
                376.297,
                0.3565,
                0.02852
        ),
        c(
                "S3 Intelligent AA",
                338.667,
                0.3565,
                0.02852
        ),
        c(
                "S3 Intelligent DAA",
                93.133,
                0.3565,
                0.02852
        ),
        c(
                "S3 Standard IA",
                1175.928,
                0.713,
                0.0713
        ),
        c(
                "S3 Express one zone",
                15051.87,
                0.17825,
                0.01426
        ),
        c(
                "S3 Glacier instant retrieval",
                376.297,
                1.426,
                0.713
        ),
        c(
                "S3 Glacier flexible retrieval",
                338.667,
                2.139,
                0.2852
        ),
        c(
                "S3 Glacier deep archive",
                93.1334,
                3.565,
                0.02852
        ),
        c(
                "S3 One Zone IA",
                940.74,
                0.713,
                0.0713
        ),
        c(
                "Microsoft Azure Premium",
                14111.1,
                0.1626,
                0.01355
        ),
        c(
                "Microsoft Azure Hot",
                1693.34,
                0.4635,
                0.03565
        ),
        c(
                "Microsoft Azure Cool",
                940.7,
                0.9269,
                0.09269
        ),
        c(
                "Microsoft Azure Cold",
                338.667,
                1.6684,
                0.92690
        ),
        c(
                "Microsoft Azure Archive",
                188.15,
                0.9269,
                46.345
        )
)

plot_cloud_costs <- function(price_list, total_size, increments, read_requests) {
  names <- sapply(pricelist, function(x) x[1])
  storage_costs <- as.numeric(sapply(pricelist, function(x) x[2]))
  read_costs <- as.numeric(sapply(pricelist, function(x) x[3]))
  upload_costs <- as.numeric(sapply(pricelist, function(x) x[4]))
  
  cost_data <- list()
  for (j in 1:length(read_costs)) {
    costs <- numeric(increments + 1)
    for (i in 0:increments) {
      costs[i + 1] <- read_costs[j] * i + storage_costs[j] * total_size + upload_costs[j] * read_requests
    }
    cost_data[[j]] <- costs
  }
  
  plot(0:increments, cost_data[[1]], type = "l", 
       xlab = "Number of Reads (i)", ylab = "Total Cost", 
       main = "Cloud Service Costs",
       ylim = range(unlist(cost_data)),
       col = 1)

  if (length(read_costs) > 1) {
    for (j in 2:length(read_costs)) {
      lines(0:increments, cost_data[[j]], col = j)
    }
  }

  legend("topleft", legend = names, col = 1:length(read_costs), lty = 1)
}

total_size <- 50
increments <- 5000
read_requests <- 9000

plot_cloud_costs(price_list, total_size, increments, read_requests)












pricelist <- list(
        c(
                "S3 Standard",
                2163.71,
                0.3565,
                0.02852
        ),
        c(
                "S3 Standard IA",
                1175.928,
                0.713,
                0.0713
        ),
        c(
                "S3 Glacier instant retrieval",
                376.297,
                1.426,
                0.713
        ),
        c(
                "S3 One Zone IA",
                940.74,
                0.713,
                0.0713
        ),
        c(
                "Microsoft Azure Hot",
                1693.34,
                0.4635,
                0.03565
        ),
        c(
                "Microsoft Azure Cool",
                940.7,
                0.9269,
                0.09269
        ),
        c(
                "Microsoft Azure Cold",
                338.667,
                1.6684,
                0.92690
        )
)
plot_cloud_costs(pricelist, total_size, increments, read_requests)