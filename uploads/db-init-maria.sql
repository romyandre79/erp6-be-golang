-- --------------------------------------------------------
-- Host:                         localhost
-- Server version:               10.11.10-MariaDB - mariadb.org binary distribution
-- Server OS:                    Win64
-- HeidiSQL Version:             12.12.0.7122
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping structure for table capella.language
DROP TABLE IF EXISTS `language`;
CREATE TABLE IF NOT EXISTS `language` (
  `languageid` int(11) NOT NULL AUTO_INCREMENT,
  `languagename` varchar(20) NOT NULL,
  `shortlang` varchar(5) NOT NULL,
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`languageid`),
  UNIQUE KEY `uq_language_name` (`languagename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.catalogsys
DROP TABLE IF EXISTS `catalogsys`;
CREATE TABLE IF NOT EXISTS `catalogsys` (
  `catalogsysid` int(11) NOT NULL AUTO_INCREMENT,
  `languageid` int(11) NOT NULL,
  `catalogname` varchar(50) NOT NULL,
  `catalogval` text NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`catalogsysid`),
  UNIQUE KEY `uq_catalogsy` (`languageid`,`catalogname`),
  KEY `fk_catalogsys_lang` (`languageid`),
  CONSTRAINT `fk_catalogsys_lang` FOREIGN KEY (`languageid`) REFERENCES `language` (`languageid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping data for table capella.catalogsys: ~0 rows (approximately)

-- Dumping structure for table capella.country
DROP TABLE IF EXISTS `country`;
CREATE TABLE IF NOT EXISTS `country` (
  `countryid` int(11) NOT NULL,
  `countrycode` varchar(2) DEFAULT NULL,
  `countryname` varchar(30) NOT NULL,
  `recordstatus` int(11) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`countryid`),
  UNIQUE KEY `uq_country_code` (`countrycode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.province
DROP TABLE IF EXISTS `province`;
CREATE TABLE IF NOT EXISTS `province` (
  `provinceid` int(11) NOT NULL AUTO_INCREMENT,
  `countryid` int(11) NOT NULL,
  `provincecode` varchar(2) DEFAULT NULL,
  `provincename` varchar(30) NOT NULL,
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`provinceid`),
  KEY `fk_province_country` (`countryid`),
  CONSTRAINT `fk_province_country` FOREIGN KEY (`countryid`) REFERENCES `country` (`countryid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.city
DROP TABLE IF EXISTS `city`;
CREATE TABLE IF NOT EXISTS `city` (
  `cityid` int(11) NOT NULL AUTO_INCREMENT,
  `provinceid` int(11) NOT NULL,
  `citycode` varchar(5) NOT NULL,
  `cityname` varchar(40) NOT NULL,
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`cityid`),
  UNIQUE KEY `uq_city_prov` (`provinceid`,`cityname`),
  CONSTRAINT `fk_city_prov` FOREIGN KEY (`provinceid`) REFERENCES `province` (`provinceid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.componentcategory
DROP TABLE IF EXISTS `componentcategory`;
CREATE TABLE IF NOT EXISTS `componentcategory` (
  `componentcategoryid` int(11) NOT NULL AUTO_INCREMENT,
  `categoryname` varchar(50) NOT NULL,
  PRIMARY KEY (`componentcategoryid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- Dumping structure for table capella.component
DROP TABLE IF EXISTS `component`;
CREATE TABLE IF NOT EXISTS `component` (
  `componentid` int(11) NOT NULL AUTO_INCREMENT,
  `componentname` varchar(50) NOT NULL,
  `componenttitle` varchar(100) NOT NULL,
  `componentcategoryid` int(11) NOT NULL,
  `componentclass` varchar(100) NOT NULL,
  `createdby` varchar(100) NOT NULL,
  `version` varchar(10) NOT NULL,
  `input` tinyint(4) NOT NULL DEFAULT 0 COMMENT 'Number of Input',
  `output` tinyint(4) NOT NULL DEFAULT 0 COMMENT 'Number of Output',
  PRIMARY KEY (`componentid`),
  UNIQUE KEY `uq_component_name` (`componentname`),
  KEY `fk_component_category` (`componentcategoryid`),
  CONSTRAINT `fk_component_category` FOREIGN KEY (`componentcategoryid`) REFERENCES `componentcategory` (`componentcategoryid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping structure for table capella.componentdetail
DROP TABLE IF EXISTS `componentdetail`;
CREATE TABLE IF NOT EXISTS `componentdetail` (
  `componentdetailid` int(11) NOT NULL AUTO_INCREMENT,
  `componentid` int(11) NOT NULL,
  `detailtype` varchar(50) DEFAULT NULL COMMENT 'Property, Event',
  `lable` varchar(50) DEFAULT NULL,
  `inputtype` varchar(100) DEFAULT NULL COMMENT 'textbox, checkbox',
  `inputname` varchar(50) DEFAULT NULL,
  `inputdesc` varchar(100) DEFAULT NULL,
  `order` smallint(6) DEFAULT 0,
  `datasourcetype` varchar(50) DEFAULT NULL COMMENT 'Table, API, List, Parameter',
  `datasource` varchar(100) DEFAULT NULL,
  `datasourceidfield` varchar(100) DEFAULT NULL,
  `datasourcenamefield` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`componentdetailid`),
  KEY `fk_comdetail_component` (`componentid`),
  CONSTRAINT `fk_comdetail_component` FOREIGN KEY (`componentid`) REFERENCES `component` (`componentid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping structure for table capella.dbobject
DROP TABLE IF EXISTS `dbobject`;
CREATE TABLE IF NOT EXISTS `dbobject` (
  `dbobjectid` int(11) NOT NULL AUTO_INCREMENT,
  `objectname` varchar(200) NOT NULL,
  `dbobjecttypeid` int(11) NOT NULL DEFAULT 0 COMMENT 'Table,StoreProcedure,Trigger,View',
  `objectcontent` text NOT NULL COMMENT 'SP, Trigger, View Content',
  `objectversion` int(11) NOT NULL DEFAULT 0,
  `ispublished` int(11) NOT NULL DEFAULT 0,
  `comment` varchar(50) DEFAULT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`dbobjectid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- Dumping structure for table capella.dbobjecttype
DROP TABLE IF EXISTS `dbobjecttype`;
CREATE TABLE IF NOT EXISTS `dbobjecttype` (
  `dbobjecttypeid` int(11) NOT NULL AUTO_INCREMENT,
  `objecttypename` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`dbobjecttypeid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- Dumping structure for table capella.groupaccess
DROP TABLE IF EXISTS `groupaccess`;
CREATE TABLE IF NOT EXISTS `groupaccess` (
  `groupaccessid` int(11) NOT NULL AUTO_INCREMENT,
  `groupname` varchar(50) NOT NULL,
  `description` varchar(100) NOT NULL,
  `recordstatus` tinyint(4) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`groupaccessid`) USING BTREE,
  UNIQUE KEY `uq_groupaccess_group` (`groupname`) USING BTREE,
  KEY `ix_groupaccess` (`groupaccessid`,`groupname`,`recordstatus`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

  -- Dumping structure for table capella.modules
DROP TABLE IF EXISTS `modules`;
CREATE TABLE IF NOT EXISTS `modules` (
  `moduleid` int(11) NOT NULL AUTO_INCREMENT,
  `modulename` varchar(30) NOT NULL,
  `description` varchar(150) NOT NULL,
  `createdby` varchar(50) NOT NULL,
  `moduleversion` varchar(3) NOT NULL DEFAULT '0',
  `installdate` timestamp NOT NULL DEFAULT current_timestamp(),
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (`moduleid`),
  UNIQUE KEY `uq_module_name` (`modulename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


  -- Dumping structure for table capella.menuaccess
DROP TABLE IF EXISTS `menuaccess`;
CREATE TABLE IF NOT EXISTS `menuaccess` (
  `menuaccessid` int(11) NOT NULL AUTO_INCREMENT,
  `menuname` varchar(20) NOT NULL,
  `menucode` varchar(50) DEFAULT NULL,
  `description` varchar(50) NOT NULL,
  `moduleid` int(11) NOT NULL,
  `parentid` int(11) DEFAULT NULL,
  `menuurl` varchar(30) DEFAULT NULL,
  `sortorder` int(11) NOT NULL DEFAULT 1,
  `menuicon` varchar(50) DEFAULT NULL,
  `menutype` varchar(50) NOT NULL DEFAULT 'LIST' COMMENT 'LIST, LIST-DETAIL, MASTER, MASTER-DETAIL, REPORT',
  `menuform` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `menuformdetail` mediumtext DEFAULT NULL,
  PRIMARY KEY (`menuaccessid`),
  UNIQUE KEY `uq_menuaccess_menuname` (`menuname`),
  KEY `fk_menuaccess_module` (`moduleid`),
  CONSTRAINT `fk_menuaccess_module` FOREIGN KEY (`moduleid`) REFERENCES `modules` (`moduleid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.groupmenu
DROP TABLE IF EXISTS `groupmenu`;
CREATE TABLE IF NOT EXISTS `groupmenu` (
  `groupmenuid` int(11) NOT NULL AUTO_INCREMENT,
  `groupaccessid` int(11) NOT NULL,
  `menuaccessid` int(11) NOT NULL,
  `isread` tinyint(4) NOT NULL DEFAULT 1,
  `iswrite` tinyint(4) NOT NULL DEFAULT 1,
  `ispost` tinyint(4) NOT NULL DEFAULT 1,
  `isreject` tinyint(4) NOT NULL DEFAULT 1,
  `isupload` tinyint(4) NOT NULL DEFAULT 1,
  `isdownload` tinyint(4) NOT NULL DEFAULT 1,
  `ispurge` tinyint(4) NOT NULL DEFAULT 1,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`groupmenuid`) USING BTREE,
  UNIQUE KEY `uq_groupmenu_gm` (`groupaccessid`,`menuaccessid`),
  KEY `fk_groupmenu_menu` (`menuaccessid`),
  CONSTRAINT `fk_groupmenu_menu` FOREIGN KEY (`menuaccessid`) REFERENCES `menuaccess` (`menuaccessid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.groupmenuauth
DROP TABLE IF EXISTS `groupmenuauth`;
CREATE TABLE IF NOT EXISTS `groupmenuauth` (
  `groupmenuauthid` int(10) NOT NULL AUTO_INCREMENT,
  `groupaccessid` int(10) NOT NULL,
  `menuauthid` int(10) NOT NULL,
  `menuvalueid` varchar(20) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`groupmenuauthid`),
  UNIQUE KEY `UQ_GROUP_MENUAUTH_VALUE` (`groupaccessid`,`menuauthid`,`menuvalueid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.menuauth
DROP TABLE IF EXISTS `menuauth`;
CREATE TABLE IF NOT EXISTS `menuauth` (
  `menuauthid` int(10) NOT NULL AUTO_INCREMENT,
  `menuobject` varchar(50) NOT NULL,
  `recordstatus` tinyint(3) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`menuauthid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.modulerelation
DROP TABLE IF EXISTS `modulerelation`;
CREATE TABLE IF NOT EXISTS `modulerelation` (
  `modulerelationid` int(11) NOT NULL AUTO_INCREMENT,
  `moduleid` int(11) NOT NULL,
  `relationid` int(11) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`modulerelationid`),
  KEY `fk_modulerel_module` (`moduleid`),
  KEY `fk_modulerel_rel` (`relationid`),
  CONSTRAINT `fk_modulerel_module` FOREIGN KEY (`moduleid`) REFERENCES `modules` (`moduleid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_modulerel_rel` FOREIGN KEY (`relationid`) REFERENCES `modules` (`moduleid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping data for table capella.modulerelation: ~0 rows (approximately)

-- Dumping structure for table capella.parameter
DROP TABLE IF EXISTS `parameter`;
CREATE TABLE IF NOT EXISTS `parameter` (
  `parameterid` int(11) NOT NULL AUTO_INCREMENT,
  `paramname` varchar(30) NOT NULL,
  `paramdata` text NOT NULL,
  `description` varchar(100) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`parameterid`) USING BTREE,
  KEY `ix_parameter` (`paramname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.snro
DROP TABLE IF EXISTS `snro`;
CREATE TABLE IF NOT EXISTS `snro` (
  `snroid` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(50) NOT NULL,
  `formatdoc` varchar(50) NOT NULL,
  `formatno` varchar(10) NOT NULL,
  `repeatby` varchar(30) DEFAULT NULL,
  `recordstatus` tinyint(4) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`snroid`),
  KEY `ix_snro` (`snroid`,`description`,`formatdoc`,`formatno`,`repeatby`,`recordstatus`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='Specific Number Range Object';

-- Dumping structure for table capella.snrodet
DROP TABLE IF EXISTS `snrodet`;
CREATE TABLE IF NOT EXISTS `snrodet` (
  `snrodetid` int(11) NOT NULL AUTO_INCREMENT,
  `plantid` int(11) DEFAULT NULL,
  `snroid` int(11) NOT NULL,
  `curdd` int(11) DEFAULT NULL,
  `curmm` int(11) DEFAULT NULL,
  `curyy` int(11) DEFAULT NULL,
  `curvalue` int(11) NOT NULL,
  `companyid` int(11) DEFAULT NULL,
  `accountid` int(11) DEFAULT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`snrodetid`) USING BTREE,
  KEY `fk_snrodet_plant` (`plantid`),
  KEY `fk_snrodet_comp` (`companyid`),
  KEY `fk_snrodet_acc` (`accountid`),
  KEY `ix_snro` (`snrodetid`,`snroid`,`curdd`,`curmm`,`curyy`,`curvalue`,`plantid`) USING BTREE
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping data for table capella.snrodet: ~0 rows (approximately)

-- Dumping structure for table capella.theme
DROP TABLE IF EXISTS `theme`;
CREATE TABLE IF NOT EXISTS `theme` (
  `themeid` int(11) NOT NULL AUTO_INCREMENT,
  `themename` varchar(20) NOT NULL,
  `description` text NOT NULL,
  `createdby` varchar(30) NOT NULL,
  `themeversion` varchar(3) NOT NULL,
  `themedata` text NOT NULL,
  `installdate` timestamp NOT NULL DEFAULT current_timestamp(),
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`themeid`),
  UNIQUE KEY `uq_theme` (`themename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.translog
DROP TABLE IF EXISTS `translog`;
CREATE TABLE IF NOT EXISTS `translog` (
  `translogid` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `createddate` timestamp NOT NULL DEFAULT current_timestamp(),
  `useraction` varchar(20) NOT NULL,
  `newdata` longtext DEFAULT NULL,
  `menuname` varchar(30) DEFAULT NULL,
  `tableid` int(11) DEFAULT NULL,
  `ippublic` varchar(100) DEFAULT NULL,
  `iplocal` varchar(100) DEFAULT NULL,
  `lat` varchar(50) DEFAULT NULL,
  `lng` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`translogid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.useraccess
DROP TABLE IF EXISTS `useraccess`;
CREATE TABLE IF NOT EXISTS `useraccess` (
  `useraccessid` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `realname` varchar(50) NOT NULL,
  `password` varchar(200) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `phoneno` varchar(20) DEFAULT NULL,
  `languageid` int(11) NOT NULL,
  `themeid` int(11) NOT NULL,
  `isonline` tinyint(4) NOT NULL DEFAULT 0,
  `joindate` datetime NOT NULL DEFAULT current_timestamp(),
  `authkey` varchar(128) DEFAULT NULL,
  `userphoto` varchar(50) DEFAULT 'man.png',
  `usertelegram` varchar(50) DEFAULT NULL,
  `userwa` varchar(50) DEFAULT NULL,
  `wallpaper` varchar(50) DEFAULT 'Andromeda-Galaxy.jpg',
  `identityid` varchar(50) DEFAULT NULL,
  `signature` varchar(50) DEFAULT NULL,
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `lastlogin` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`useraccessid`),
  UNIQUE KEY `uq_username` (`username`),
  UNIQUE KEY `uq_useremail` (`email`),
  KEY `fk_user_lang` (`languageid`),
  KEY `fk_user_theme` (`themeid`),
  CONSTRAINT `fk_user_lang` FOREIGN KEY (`languageid`) REFERENCES `language` (`languageid`) ON UPDATE CASCADE,
  CONSTRAINT `fk_user_theme` FOREIGN KEY (`themeid`) REFERENCES `theme` (`themeid`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.widget
DROP TABLE IF EXISTS `widget`;
CREATE TABLE IF NOT EXISTS `widget` (
  `widgetid` int(11) NOT NULL AUTO_INCREMENT,
  `widgetname` varchar(50) NOT NULL,
  `widgettitle` varchar(40) NOT NULL,
  `widgetversion` varchar(3) NOT NULL DEFAULT '0',
  `widgetby` varchar(20) NOT NULL,
  `description` varchar(100) NOT NULL,
  `widgetform` text NOT NULL,
  `moduleid` int(11) DEFAULT NULL,
  `installdate` timestamp NOT NULL DEFAULT current_timestamp(),
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`widgetid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.userdash
DROP TABLE IF EXISTS `userdash`;
CREATE TABLE IF NOT EXISTS `userdash` (
  `userdashid` int(11) NOT NULL AUTO_INCREMENT,
  `groupaccessid` int(11) NOT NULL,
  `widgetid` int(11) NOT NULL,
  `menuaccessid` int(11) NOT NULL,
  `position` tinyint(4) NOT NULL DEFAULT 0,
  `webformat` varchar(20) NOT NULL DEFAULT '0',
  `dashgroup` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`userdashid`),
  KEY `fk_userwidget_group` (`groupaccessid`),
  KEY `fk_userwdget_menu` (`menuaccessid`),
  KEY `fk_userwidge_widget` (`widgetid`),
  CONSTRAINT `fk_userwdget_menu` FOREIGN KEY (`menuaccessid`) REFERENCES `menuaccess` (`menuaccessid`),
  CONSTRAINT `fk_userwidge_widget` FOREIGN KEY (`widgetid`) REFERENCES `widget` (`widgetid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.userfav
DROP TABLE IF EXISTS `userfav`;
CREATE TABLE IF NOT EXISTS `userfav` (
  `userfavid` int(11) NOT NULL AUTO_INCREMENT,
  `useraccessid` int(11) NOT NULL,
  `menuaccessid` int(11) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`userfavid`),
  KEY `fk_userfav_menu` (`menuaccessid`),
  KEY `fk_userfav_user` (`useraccessid`),
  CONSTRAINT `fk_userfav_menu` FOREIGN KEY (`menuaccessid`) REFERENCES `menuaccess` (`menuaccessid`),
  CONSTRAINT `fk_userfav_user` FOREIGN KEY (`useraccessid`) REFERENCES `useraccess` (`useraccessid`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping data for table capella.userfav: ~0 rows (approximately)

-- Dumping structure for table capella.usergroup
DROP TABLE IF EXISTS `usergroup`;
CREATE TABLE IF NOT EXISTS `usergroup` (
  `usergroupid` int(11) NOT NULL AUTO_INCREMENT,
  `useraccessid` int(11) NOT NULL,
  `groupaccessid` int(11) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`usergroupid`),
  UNIQUE KEY `uq_groupmenu` (`groupaccessid`,`useraccessid`),
  KEY `fk_usergroup_user` (`useraccessid`),
  KEY `fk_usergroup_group` (`groupaccessid`),
  CONSTRAINT `fk_usergroup_group` FOREIGN KEY (`groupaccessid`) REFERENCES `groupaccess` (`groupaccessid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.userinbox
DROP TABLE IF EXISTS `userinbox`;
CREATE TABLE IF NOT EXISTS `userinbox` (
  `userinboxid` int(11) NOT NULL AUTO_INCREMENT,
  `useraccessid` int(11) NOT NULL,
  `fromuserid` int(11) NOT NULL,
  `subject` varchar(150) NOT NULL,
  `description` text NOT NULL,
  `senddate` timestamp NOT NULL DEFAULT current_timestamp(),
  `recordstatus` tinyint(4) NOT NULL DEFAULT 1,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`userinboxid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping data for table capella.userinbox: ~0 rows (approximately)

-- Dumping structure for table capella.usertodo
DROP TABLE IF EXISTS `usertodo`;
CREATE TABLE IF NOT EXISTS `usertodo` (
  `usertodoid` int(11) NOT NULL AUTO_INCREMENT,
  `useraccessid` int(11) NOT NULL,
  `tododate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `menuname` varchar(50) DEFAULT NULL,
  `docno` varchar(200) DEFAULT NULL,
  `description` varchar(400) DEFAULT NULL,
  `isread` tinyint(4) DEFAULT 0,
  `recordstatus` tinyint(4) NOT NULL DEFAULT 0,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`usertodoid`),
  KEY `fk_usertodo_user` (`useraccessid`),
  CONSTRAINT `fk_usertodo_user` FOREIGN KEY (`useraccessid`) REFERENCES `useraccess` (`useraccessid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping data for table capella.usertodo: ~0 rows (approximately)

-- Dumping structure for table capella.wfcomponent
DROP TABLE IF EXISTS `wfcomponent`;
CREATE TABLE IF NOT EXISTS `wfcomponent` (
  `wfcomponentid` int(11) NOT NULL AUTO_INCREMENT,
  `workflowid` int(11) NOT NULL,
  `componentname` varchar(50) NOT NULL,
  `componentid` int(11) NOT NULL,
  `componentdetailid` int(11) NOT NULL,
  `wfvalue` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`wfcomponentid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table capella.wfcomponent: ~0 rows (approximately)

-- Dumping structure for table capella.wfgroup
DROP TABLE IF EXISTS `wfgroup`;
CREATE TABLE IF NOT EXISTS `wfgroup` (
  `wfgroupid` int(11) NOT NULL AUTO_INCREMENT,
  `workflowid` int(11) NOT NULL,
  `groupaccessid` int(11) NOT NULL,
  `wfbefstat` tinyint(4) NOT NULL,
  `wfrecstat` tinyint(4) NOT NULL,
  `sla` int(11) NOT NULL DEFAULT 1,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`wfgroupid`) USING BTREE,
  UNIQUE KEY `ix_wfgroup_wgb` (`workflowid`,`groupaccessid`,`wfbefstat`) USING BTREE,
  KEY `ix_wfgroup_wfgbr` (`workflowid`,`groupaccessid`,`wfbefstat`,`wfrecstat`),
  KEY `ix_wfgroup_wgr` (`workflowid`,`groupaccessid`,`wfrecstat`),
  KEY `fk_wfgroup_group` (`groupaccessid`),
  CONSTRAINT `fk_wfgroup_group` FOREIGN KEY (`groupaccessid`) REFERENCES `groupaccess` (`groupaccessid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_wfgroup_wf` FOREIGN KEY (`workflowid`) REFERENCES `workflow` (`workflowid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.wfstatus
DROP TABLE IF EXISTS `wfstatus`;
CREATE TABLE IF NOT EXISTS `wfstatus` (
  `wfstatusid` int(11) NOT NULL AUTO_INCREMENT,
  `workflowid` int(11) NOT NULL,
  `wfstat` tinyint(4) NOT NULL,
  `wfstatusname` varchar(50) NOT NULL,
  `backcolor` varchar(50) NOT NULL DEFAULT 'red',
  `fontcolor` varchar(50) NOT NULL DEFAULT 'white',
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`wfstatusid`),
  UNIQUE KEY `uq_wfstatus` (`workflowid`,`wfstat`),
  KEY `ix_wfstatus` (`wfstatusid`,`workflowid`,`wfstat`,`wfstatusname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


-- Dumping structure for table capella.workflow
DROP TABLE IF EXISTS `workflow`;
CREATE TABLE IF NOT EXISTS `workflow` (
  `workflowid` int(11) NOT NULL AUTO_INCREMENT,
  `wfname` varchar(50) NOT NULL,
  `wfdesc` varchar(200) NOT NULL COMMENT 'wf description',
  `wfminstat` tinyint(4) NOT NULL,
  `wfmaxstat` tinyint(4) NOT NULL,
  `flow` mediumtext DEFAULT NULL,
  `recordstatus` tinyint(4) NOT NULL,
  `updatedate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`workflowid`),
  UNIQUE KEY `uq_workflow_wfname` (`wfname`),
  KEY `ix_workflow_wfname` (`wfname`),
  KEY `ix_workflow_id` (`workflowid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Dumping structure for table capella.workflowdetail
DROP TABLE IF EXISTS `workflowdetail`;
CREATE TABLE IF NOT EXISTS `workflowdetail` (
  `workflowdetailid` bigint(20) NOT NULL AUTO_INCREMENT,
  `componentid` int(11) NOT NULL,
  `componentdetailid` int(11) NOT NULL,
  `workflowid` int(11) NOT NULL,
  `componentvalue` mediumtext NOT NULL,
  `nodeid` int(11) NOT NULL,
  PRIMARY KEY (`workflowdetailid`),
  KEY `fk_wfdetail_component` (`componentid`),
  KEY `fk_wfdetail_comdetail` (`componentdetailid`),
  KEY `fk_wfdetail_wfid` (`workflowid`),
  CONSTRAINT `fk_wfdetail_comdetail` FOREIGN KEY (`componentdetailid`) REFERENCES `componentdetail` (`componentdetailid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_wfdetail_component` FOREIGN KEY (`componentid`) REFERENCES `component` (`componentid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_wfdetail_wfid` FOREIGN KEY (`workflowid`) REFERENCES `workflow` (`workflowid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- Dumping structure for table capella.workflowparameter
DROP TABLE IF EXISTS `workflowparameter`;
CREATE TABLE IF NOT EXISTS `workflowparameter` (
  `wfparameterid` int(11) NOT NULL AUTO_INCREMENT,
  `workflowid` int(11) NOT NULL,
  `parametername` varchar(100) NOT NULL,
  `parametervalue` varchar(100) NOT NULL,
  `parametertype` varchar(50) NOT NULL,
  PRIMARY KEY (`wfparameterid`),
  UNIQUE KEY `uq_wfparameter` (`workflowid`,`parametername`),
  CONSTRAINT `fk_wfparameter_workflow` FOREIGN KEY (`workflowid`) REFERENCES `workflow` (`workflowid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- Dumping structure for procedure capella.InsertTransLogEx
DROP PROCEDURE IF EXISTS `InsertTransLogEx`;
CREATE PROCEDURE `InsertTransLogEx`(
	IN `vuseraction` VARCHAR(50),
	IN `vnewdata` TEXT,
	IN `vmenuname` VARCHAR(50),
	IN `vtableid` INT,
	IN `vdatauser` TEXT
)
BEGIN
	declare vlog int;
	declare vippublic,viplocal,vlat,vlng,vusername text;
		
	set vusername = Split_Str(vdatauser, ',', 1);
	set vippublic = Split_Str(vdatauser, ',', 2);
	set viplocal = Split_Str(vdatauser, ',', 3);
	set vlat = Split_Str(vdatauser, ',', 4);
	set vlng = Split_Str(vdatauser, ',', 5);
  
  SELECT a.paramdata
	into vlog
	from parameter a 
	where UPPER(a.paramname) = 'USINGLOG'
  LIMIT 1;


	if vlog = 1 then
		insert into translog (username,useraction,newdata,menuname,tableid,ippublic,iplocal,lat,lng)
    value (vusername,vuseraction,vnewdata,vmenuname,vtableid,vippublic,viplocal,vlat,vlng);
	end if;
END;

-- Dumping structure for procedure capella.ModifWorkflow
DROP PROCEDURE IF EXISTS `ModifWorkflow`;
CREATE PROCEDURE `ModifWorkflow`(
	IN `vid` INT,
	IN `vwfname` VARCHAR(50),
	IN `vwfdesc` VARCHAR(100),
	IN `vwfminstat` integer,
	IN `vwfmaxstat` integer,
	IN `vrecordstatus` INT,
	IN `vdatauser` TEXT
)
BEGIN
	declare k int;
	
	if (vid < 0) then
		insert into workflow (wfname,wfdesc,wfminstat,wfmaxstat,recordstatus)
		values (vwfname,vwfdesc,vwfminstat,vwfmaxstat,vrecordstatus);
		set k = last_insert_id();
		call inserttranslogex('new',
			concat('wfname=',coalesce(vwfname,''),
				',wfdesc=',coalesce(vwfdesc,''),
				',wfminstat=',coalesce(vwfminstat,''),
				',wfmaxstat=',coalesce(vwfmaxstat,''),
				',recordstatus=',coalesce(vrecordstatus,'')),
			'workflow',k,vdatauser);
	else
		update workflow
		set wfname = vwfname, wfdesc = vwfdesc, wfminstat = vwfminstat,
	    wfmaxstat = vwfmaxstat, recordstatus=vrecordstatus
		where workflowid = vid;
		set k = vid;
		call inserttranslogex('update',
			concat('wfname=',coalesce(vwfname,''),
				',wfdesc=',coalesce(vwfdesc,''),
	      	',wfminstat=',coalesce(vwfminstat,''),
				',wfmaxstat=',coalesce(vwfmaxstat,''),
				',recordstatus=',coalesce(vrecordstatus,'')),'workflow',vid,vdatauser);
	end if;
	update wfgroup
	set workflowid = k
	where workflowid = vid;
	update wfstatus
	set workflowid = k
	where workflowid = vid;
END;

-- Dumping structure for procedure capella.UpdateMenuFlow
DROP PROCEDURE IF EXISTS `UpdateMenuFlow`;
CREATE PROCEDURE `UpdateMenuFlow`(
	IN `vworkflowid` INT,
	IN `vmenuflow` LONGTEXT,
	IN `vdatauser` TEXT
)
BEGIN
	call inserttranslogex('update',
		concat('flow=',coalesce(vmenuflow,''))
  ,'workflow',vworkflowid,vdatauser);

	update workflow 
	set flow = vmenuflow
	where workflowid=vworkflowid;
END;

-- Dumping structure for function capella.Split_Str
DROP FUNCTION IF EXISTS `Split_Str`;
CREATE FUNCTION `Split_Str`(`x` VARCHAR(255), `delim` VARCHAR(12), `pos` INT
				) RETURNS varchar(255) CHARSET utf8mb3 COLLATE utf8mb3_general_ci
BEGIN
				RETURN REPLACE(SUBSTRING(SUBSTRING_INDEX(x, delim, pos),
							 LENGTH(SUBSTRING_INDEX(x, delim, pos -1)) + 1),
							 delim, '');
				END

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
