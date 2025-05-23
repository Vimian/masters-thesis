CREATE TABLE IF NOT EXISTS measurements (
        id SERIAL,
        algorithm VARCHAR(255),
        run INT,

        original_path VARCHAR(255),
        original_size BIGINT,
        original_hash VARCHAR(255),
        PRIMARY KEY (algorithm, run, original_path),

        compress_start_get_object_info BIGINT,
        compress_end_get_object_info BIGINT,
        compress_duration_get_object_info BIGINT,
        compress_start_get_object BIGINT,
        compress_end_get_object BIGINT,
        compress_duration_get_object BIGINT,
        compress_start_algorithm BIGINT,
        compress_end_algorithm BIGINT,
        compress_duration_algorithm BIGINT,
        compress_start_upload BIGINT,
        compress_end_upload BIGINT,
        compress_duration_upload BIGINT,
        compress_path VARCHAR(255),
        compress_size BIGINT,
        compress_hash VARCHAR(255),

        decompress_start_get_object_info BIGINT,
        decompress_end_get_object_info BIGINT,
        decompress_duration_get_object_info BIGINT,
        decompress_start_get_object BIGINT,
        decompress_end_get_object BIGINT,
        decompress_duration_get_object BIGINT,
        decompress_start_algorithm BIGINT,
        decompress_end_algorithm BIGINT,
        decompress_duration_algorithm BIGINT,
        decompress_start_upload BIGINT,
        decompress_end_upload BIGINT,
        decompress_duration_upload BIGINT,
        decompress_path VARCHAR(255),
        decompress_size BIGINT,
        decompress_hash VARCHAR(255),

        compression_ratio FLOAT
);

CREATE TABLE IF NOT EXISTS analytics (
        id SERIAL PRIMARY KEY,
        file_path VARCHAR(255),
        file_name VARCHAR(255),
        file_size BIGINT,
        window_length_bytes BIGINT,
        dictionary_length BIGINT,
        dictionary_limit_reached SMALLINT,
        compressed_size_bytes_go BIGINT
);

CREATE TABLE IF NOT EXISTS cloud (
        id SERIAL PRIMARY KEY,
        tier_name VARCHAR(255),
        run INT,

        size BIGINT,
        original_hash VARCHAR(255),
        downloaded_hash VARCHAR(255),

        start_upload BIGINT,
        end_upload BIGINT,
        duration_upload BIGINT,
        start_download BIGINT,
        end_download BIGINT,
        duration_download BIGINT
);