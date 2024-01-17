package app

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const (
	DataSourceName = "mypass.sqlite"
)

type DBStruct struct {
	db   *sql.DB
	once sync.Once
}

var Sqlite DBStruct

func (s *DBStruct) DB() *sql.DB {
	s.once.Do(func() {
		var err error
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			Log.Info(fmt.Sprintf("Could not find local user folder. Error: %v\n", err))
		}
		s.db, err = sql.Open("sqlite3", userHomeDir+string(os.PathSeparator)+DataSourceName)
		if err != nil {
			panic(err.Error())
		}
	})
	return s.db
}

func BackupDatabase() {
	const maxBackups = 15
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		Log.Error(fmt.Sprintf("Could not find local user folder. Error: %v\n", err))
	}
	backupDir := userHomeDir + string(os.PathSeparator) + ".mypass_backups"
	dbFile := userHomeDir + string(os.PathSeparator) + DataSourceName
	// 创建备份目录（如果不存在）
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.Mkdir(backupDir, os.ModePerm)
	}
	// 获取最新备份的 MD5
	latestBackup, latestMD5, err := getLatestBackupMD5(backupDir)
	if err != nil {
		Log.Error(fmt.Sprintf("Error getting latest backup: %v\n", err))
		return
	}
	Log.Info(latestBackup)
	// 获取当前数据库文件的 MD5
	currentMD5, err := getFileMD5(dbFile)
	if err != nil {
		Log.Error(fmt.Sprintf("Error getting current DB MD5: %v\n", err))
		return
	}

	// 比较 MD5，如果一致，则不备份
	if latestMD5 == currentMD5 {
		Log.Info("No changes in database, skipping backup.")
		return
	}

	// 创建新的备份文件
	backupFileName := fmt.Sprintf("%s/backup_%s.db", backupDir, time.Now().Format("2006-01-02_15-04-05"))
	err = copyFile(dbFile, backupFileName)
	if err != nil {
		Log.Error(fmt.Sprintf("Error creating backup: %v\n", err))
		return
	}

	// 保留最新的 15 个备份，删除其余的
	err = pruneOldBackups(backupDir, maxBackups)
	if err != nil {
		Log.Error(fmt.Sprintf("Error pruning old backups: %v\n", err))
		return
	}
	Log.Info(fmt.Sprintf("Backup created: %s\n", backupFileName))
}

// getFileMD5 计算文件的 MD5 值
func getFileMD5(filePath string) (string, error) {
	var md5String string
	file, err := os.Open(filePath)
	if err != nil {
		return md5String, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	md5String = fmt.Sprintf("%x", hashInBytes)
	return md5String, nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// getLatestBackupMD5 返回最新备份文件的路径和其 MD5 值
func getLatestBackupMD5(backupDir string) (string, string, error) {
	dirEntries, _ := os.ReadDir(backupDir)
	sort.Slice(dirEntries, func(i, j int) bool {
		ifileInfo, _ := dirEntries[i].Info()
		jfileInfo, _ := dirEntries[j].Info()
		return ifileInfo.ModTime().After(jfileInfo.ModTime())
	})

	for _, f := range dirEntries {
		if filepath.Ext(f.Name()) == ".db" {
			fullPath := filepath.Join(backupDir, f.Name())
			md5, err := getFileMD5(fullPath)
			if err != nil {
				return "", "", err
			}
			return fullPath, md5, nil
		}
	}
	return "", "", nil // 无备份文件
}

// pruneOldBackups 删除超出最大数量的旧备份文件
func pruneOldBackups(backupDir string, maxBackups int) error {
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	if len(files) > maxBackups {
		for _, file := range files[:len(files)-maxBackups] {
			err := os.Remove(filepath.Join(backupDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
