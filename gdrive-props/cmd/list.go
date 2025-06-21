package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ToshihitoKon/agent-works/gdrive-props/pkg/auth"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [folder-id]",
	Short: "List files in Google Drive",
	Long:  "List files in Google Drive. If no folder ID is specified, shows My Drive contents.",
	RunE:  runList,
}

var (
	listPageSize int64
	listShowAll  bool
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Int64VarP(&listPageSize, "page-size", "p", 50, "Number of files to display per page")
	listCmd.Flags().BoolVarP(&listShowAll, "all", "a", false, "Show all files including trashed")
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	auth, err := auth.NewAuth()
	if err != nil {
		return fmt.Errorf("failed to initialize auth: %v", err)
	}

	service, err := auth.GetDriveService(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Drive service: %v", err)
	}

	// Determine folder to list
	var folderID string
	if len(args) > 0 {
		folderID = args[0]
	}

	// Build query
	query := ""
	if folderID != "" {
		query = fmt.Sprintf("'%s' in parents", folderID)
	} else {
		// Show My Drive contents (root folder)
		query = "'root' in parents"
	}

	if !listShowAll {
		if query != "" {
			query += " and "
		}
		query += "trashed = false"
	}

	// List files
	filesList, err := service.Files.List().
		Q(query).
		PageSize(listPageSize).
		Fields("files(id,name,mimeType,size,modifiedTime,appProperties)").
		OrderBy("name").
		Do()
	if err != nil {
		return fmt.Errorf("failed to list files: %v", err)
	}

	if len(filesList.Files) == 0 {
		if folderID != "" {
			fmt.Printf("No files found in folder %s\n", folderID)
		} else {
			fmt.Println("No files found in My Drive")
		}
		return nil
	}

	// Display header
	if folderID != "" {
		fmt.Printf("Files in folder %s (%d):\n", folderID, len(filesList.Files))
	} else {
		fmt.Printf("My Drive (%d files):\n", len(filesList.Files))
	}
	
	fmt.Println()
	fmt.Printf("%-20s %-30s %-15s %-12s %-20s %s\n", 
		"ID", "Name", "Type", "Size", "Modified", "Properties")
	fmt.Println(strings.Repeat("-", 120))

	for _, file := range filesList.Files {
		// Format size
		sizeStr := "-"
		if file.Size > 0 {
			sizeStr = formatSize(file.Size)
		}

		// Format modified time
		modifiedStr := "-"
		if file.ModifiedTime != "" {
			if t, err := time.Parse(time.RFC3339, file.ModifiedTime); err == nil {
				modifiedStr = t.Format("2006-01-02 15:04")
			}
		}

		// Format file type
		typeStr := "file"
		if file.MimeType == "application/vnd.google-apps.folder" {
			typeStr = "folder"
		} else if strings.HasPrefix(file.MimeType, "application/vnd.google-apps.") {
			typeStr = "google-" + strings.TrimPrefix(file.MimeType, "application/vnd.google-apps.")
		}

		// Format properties
		propsStr := ""
		if file.AppProperties != nil && len(file.AppProperties) > 0 {
			var propPairs []string
			for k, v := range file.AppProperties {
				propPairs = append(propPairs, fmt.Sprintf("%s=%s", k, v))
			}
			propsStr = fmt.Sprintf("{%s}", strings.Join(propPairs, ", "))
		}

		fmt.Printf("%-20s %-30s %-15s %-12s %-20s %s\n",
			truncateString(file.Id, 18),
			truncateString(file.Name, 28),
			truncateString(typeStr, 13),
			sizeStr,
			modifiedStr,
			truncateString(propsStr, 30))
	}

	fmt.Printf("\nTotal: %d files\n", len(filesList.Files))
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}