/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2025/12/4 18:22:23                           */
/*==============================================================*/
create database volunteer;
use volunteer;
DROP TABLE IF EXISTS ApplicationStatusLog;
DROP TABLE IF EXISTS Application;
DROP TABLE IF EXISTS Activity;
DROP TABLE IF EXISTS ActivityCategory;
DROP TABLE IF EXISTS Dept;
DROP TABLE IF EXISTS User;
DROP TABLE IF EXISTS Role;

/*==============================================================*/
/* Table: Role                                                  */
/*==============================================================*/
CREATE TABLE Role
(
   role_id              INT NOT NULL AUTO_INCREMENT,
   role_name            VARCHAR(20) NOT NULL,
   PRIMARY KEY (role_id)
);

/*==============================================================*/
/* Table: User                                                  */
/*==============================================================*/
CREATE TABLE User
(
   user_id              INT NOT NULL AUTO_INCREMENT,
   role_id              INT NOT NULL,
   username             VARCHAR(50) NOT NULL,
   password             VARCHAR(100) NOT NULL,
   PRIMARY KEY (user_id),
   UNIQUE KEY uk_username (username)
);

/*==============================================================*/
/* Table: Dept                                                  */
/*==============================================================*/
CREATE TABLE Dept
(
   dept_id              INT NOT NULL AUTO_INCREMENT,
   dept_name            VARCHAR(50) NOT NULL,
   PRIMARY KEY (dept_id)
);

/*==============================================================*/
/* Table: ActivityCategory                                      */
/*==============================================================*/
CREATE TABLE ActivityCategory
(
   category_id          INT NOT NULL AUTO_INCREMENT,
   category_name        VARCHAR(50) NOT NULL,
   PRIMARY KEY (category_id)
);

/*==============================================================*/
/* Table: Activity                                              */
/*==============================================================*/
CREATE TABLE Activity
(
   activity_id          INT NOT NULL AUTO_INCREMENT,
   dept_id              INT NOT NULL,
   category_id          INT NOT NULL,
   creator_id           INT NOT NULL,  -- 改名为creator_id更明确
   title                VARCHAR(100) NOT NULL,
   activity_time        DATETIME NOT NULL,
   location             VARCHAR(100) NOT NULL,
   max_people           INT NOT NULL,
   PRIMARY KEY (activity_id)
);

/*==============================================================*/
/* Table: Application                                           */
/*==============================================================*/
CREATE TABLE Application
(
   application_id       INT NOT NULL AUTO_INCREMENT,
   user_id              INT NOT NULL,
   activity_id          INT NOT NULL,
   apply_time           DATETIME NOT NULL,
   current_status       VARCHAR(20) NOT NULL DEFAULT 'pending',
   PRIMARY KEY (application_id),
   UNIQUE KEY uk_user_activity (user_id, activity_id)  -- 防止重复申请
);

/*==============================================================*/
/* Table: ApplicationStatusLog                                  */
/*==============================================================*/
CREATE TABLE ApplicationStatusLog
(
   log_id               INT NOT NULL AUTO_INCREMENT,
   application_id       INT NOT NULL,
   handler_id           INT NULL,  -- 改名为handler_id更明确，允许NULL
   log_status           VARCHAR(20) NOT NULL,
   handle_time          DATETIME NOT NULL,
   PRIMARY KEY (log_id)
);

/*==============================================================*/
/* Foreign Key Constraints                                      */
/*==============================================================*/

-- User表外键
ALTER TABLE User ADD CONSTRAINT fk_user_role 
    FOREIGN KEY (role_id) REFERENCES Role (role_id);

-- Activity表外键
ALTER TABLE Activity ADD CONSTRAINT fk_activity_dept 
    FOREIGN KEY (dept_id) REFERENCES Dept (dept_id);
ALTER TABLE Activity ADD CONSTRAINT fk_activity_category 
    FOREIGN KEY (category_id) REFERENCES ActivityCategory (category_id);
ALTER TABLE Activity ADD CONSTRAINT fk_activity_creator 
    FOREIGN KEY (creator_id) REFERENCES User (user_id);

-- Application表外键
ALTER TABLE Application ADD CONSTRAINT fk_application_user 
    FOREIGN KEY (user_id) REFERENCES User (user_id);
ALTER TABLE Application ADD CONSTRAINT fk_application_activity 
    FOREIGN KEY (activity_id) REFERENCES Activity (activity_id);

-- ApplicationStatusLog表外键
ALTER TABLE ApplicationStatusLog ADD CONSTRAINT fk_statuslog_application 
    FOREIGN KEY (application_id) REFERENCES Application (application_id);
ALTER TABLE ApplicationStatusLog ADD CONSTRAINT fk_statuslog_handler 
    FOREIGN KEY (handler_id) REFERENCES User (user_id);