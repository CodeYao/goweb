-- MySQL dump 10.13  Distrib 5.7.17, for Win64 (x86_64)
--
-- Host: localhost    Database: tjfoc_ca
-- ------------------------------------------------------
-- Server version	8.0.11

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `account`
--

DROP TABLE IF EXISTS `account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `accountId` varchar(45) DEFAULT NULL,
  `organization` varchar(45) DEFAULT NULL,
  `password` varchar(45) DEFAULT NULL,
  `accountLevel` varchar(45) DEFAULT NULL,
  `enabled` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `account`
--

LOCK TABLES `account` WRITE;
/*!40000 ALTER TABLE `account` DISABLE KEYS */;
INSERT INTO `account` VALUES (1,'admin','tjfoc','e10adc3949ba59abbe56e057f20f883e','99','enabled'),(2,'007','宝武','e10adc3949ba59abbe56e057f20f883e','11','enabled');
/*!40000 ALTER TABLE `account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `address`
--

DROP TABLE IF EXISTS `address`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(256) DEFAULT NULL,
  `addDate` varchar(45) DEFAULT NULL,
  `addPeople` varchar(45) DEFAULT NULL,
  `descript` varchar(256) DEFAULT NULL,
  `enabled` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `address`
--

LOCK TABLES `address` WRITE;
/*!40000 ALTER TABLE `address` DISABLE KEYS */;
/*!40000 ALTER TABLE `address` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ca`
--

DROP TABLE IF EXISTS `ca`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ca` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `capubkey` varchar(512) DEFAULT NULL,
  `caprivkey` varchar(512) DEFAULT NULL,
  `cacert` varchar(2048) DEFAULT NULL,
  `keytype` varchar(45) DEFAULT NULL,
  `createdate` varchar(45) DEFAULT NULL,
  `accountId` varchar(45) DEFAULT NULL,
  `enabled` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ca`
--

LOCK TABLES `ca` WRITE;
/*!40000 ALTER TABLE `ca` DISABLE KEYS */;
/*!40000 ALTER TABLE `ca` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cert`
--

DROP TABLE IF EXISTS `cert`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cert` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `certname` varchar(45) DEFAULT NULL,
  `ipstr` varchar(45) DEFAULT NULL,
  `cert` varchar(1024) DEFAULT NULL,
  `pubkey` varchar(512) DEFAULT NULL,
  `accountId` varchar(45) DEFAULT NULL,
  `reqdata` varchar(1024) DEFAULT NULL,
  `reqdate` varchar(45) DEFAULT NULL,
  `state` varchar(45) DEFAULT NULL,
  `approver` varchar(45) DEFAULT NULL,
  `approvaldate` varchar(45) DEFAULT NULL,
  `revoker` varchar(45) DEFAULT NULL,
  `revokedate` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cert`
--

LOCK TABLES `cert` WRITE;
/*!40000 ALTER TABLE `cert` DISABLE KEYS */;
/*!40000 ALTER TABLE `cert` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `codeaddress`
--

DROP TABLE IF EXISTS `codeaddress`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `codeaddress` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(512) DEFAULT NULL,
  `addDate` varchar(45) DEFAULT NULL,
  `addPeople` varchar(45) DEFAULT NULL,
  `descript` varchar(256) DEFAULT NULL,
  `enabled` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `codeaddress`
--

LOCK TABLES `codeaddress` WRITE;
/*!40000 ALTER TABLE `codeaddress` DISABLE KEYS */;
/*!40000 ALTER TABLE `codeaddress` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-08-08  9:23:14
