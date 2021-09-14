-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Sep 13, 2021 at 08:04 PM
-- Server version: 8.0.26-0ubuntu0.20.04.2
-- PHP Version: 7.4.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `fusion`
--

-- --------------------------------------------------------

--
-- Table structure for table `album`
--

CREATE TABLE `album` (
  `id` varchar(128) NOT NULL,
  `name` varchar(128) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `album`
--

INSERT INTO `album` (`id`, `name`, `created_at`, `updated_at`) VALUES
('2367710a-d4fb-49f5-8860-557b337386dd', 'KIRK', '2019-10-05 02:21:11', '2019-10-05 02:21:11'),
('967d5bb5-3a7a-4d5e-8a6c-febc8c5b3f13', 'Hollywood\'s Bleeding', '2019-10-01 12:36:38', '2019-10-01 12:36:38'),
('b0a24f12-428f-4ff5-84d5-bc1fdcff6f03', 'Lover', '2019-10-11 16:43:18', '2019-10-11 16:43:18'),
('c809bf15-bc2c-4621-bb96-70af96fd5d67', 'AI YoungBoy 2', '2019-10-02 08:16:12', '2019-10-02 08:16:12'),
('e0bb80ec-75a6-4348-bfc3-6ac1e89b195e', 'So Much Fun', '2019-10-12 09:16:02', '2019-10-12 09:16:02');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `album`
--
ALTER TABLE `album`
  ADD PRIMARY KEY (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
