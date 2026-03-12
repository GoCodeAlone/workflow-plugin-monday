package internal

import sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"

type stepConstructor func(name string, config map[string]any) (sdk.StepInstance, error)

func stepRegistry() map[string]stepConstructor {
	return map[string]stepConstructor{
		// boards
		"step.monday_create_board":    newCreateBoardStep,
		"step.monday_list_boards":     newListBoardsStep,
		"step.monday_fetch_board":     newFetchBoardStep,
		"step.monday_update_board":    newUpdateBoardStep,
		"step.monday_delete_board":    newDeleteBoardStep,
		"step.monday_duplicate_board": newDuplicateBoardStep,
		"step.monday_archive_board":   newArchiveBoardStep,
		// items
		"step.monday_create_item":  newCreateItemStep,
		"step.monday_list_items":   newListItemsStep,
		"step.monday_fetch_item":   newFetchItemStep,
		"step.monday_update_item":  newUpdateItemStep,
		"step.monday_move_item":    newMoveItemStep,
		"step.monday_archive_item": newArchiveItemStep,
		"step.monday_delete_item":  newDeleteItemStep,
		"step.monday_search_items": newSearchItemsStep,
		// subitems
		"step.monday_create_subitem": newCreateSubitemStep,
		"step.monday_list_subitems":  newListSubitemsStep,
		"step.monday_update_subitem": newUpdateSubitemStep,
		"step.monday_delete_subitem": newDeleteSubitemStep,
		// columns
		"step.monday_get_column_values":   newGetColumnValuesStep,
		"step.monday_change_column_value": newChangeColumnValueStep,
		"step.monday_create_column":       newCreateColumnStep,
		// groups
		"step.monday_create_group": newCreateGroupStep,
		"step.monday_list_groups":  newListGroupsStep,
		"step.monday_update_group": newUpdateGroupStep,
		"step.monday_move_group":   newMoveGroupStep,
		"step.monday_delete_group": newDeleteGroupStep,
		// workspaces
		"step.monday_create_workspace": newCreateWorkspaceStep,
		"step.monday_list_workspaces":  newListWorkspacesStep,
		"step.monday_update_workspace": newUpdateWorkspaceStep,
		"step.monday_delete_workspace": newDeleteWorkspaceStep,
		// folders
		"step.monday_create_folder": newCreateFolderStep,
		"step.monday_list_folders":  newListFoldersStep,
		"step.monday_update_folder": newUpdateFolderStep,
		"step.monday_delete_folder": newDeleteFolderStep,
		// updates
		"step.monday_create_update": newCreateUpdateStep,
		"step.monday_list_updates":  newListUpdatesStep,
		"step.monday_edit_update":   newEditUpdateStep,
		"step.monday_delete_update": newDeleteUpdateStep,
		// users
		"step.monday_list_users":  newListUsersStep,
		"step.monday_fetch_user":  newFetchUserStep,
		"step.monday_invite_user": newInviteUserStep,
		// teams
		"step.monday_list_teams":             newListTeamsStep,
		"step.monday_add_team_to_workspace":  newAddTeamToWorkspaceStep,
		// tags
		"step.monday_list_tags":  newListTagsStep,
		"step.monday_create_tag": newCreateTagStep,
		// files
		"step.monday_upload_file": newUploadFileStep,
		"step.monday_list_files":  newListFilesStep,
		// notifications
		"step.monday_create_notification": newCreateNotificationStep,
		// webhooks
		"step.monday_create_webhook": newCreateWebhookStep,
		"step.monday_list_webhooks":  newListWebhooksStep,
		"step.monday_delete_webhook": newDeleteWebhookStep,
		// documents
		"step.monday_create_document": newCreateDocumentStep,
		"step.monday_list_documents":  newListDocumentsStep,
		"step.monday_update_document": newUpdateDocumentStep,
		// generic
		"step.monday_query":  newQueryStep,
		"step.monday_mutate": newMutateStep,
	}
}

func allStepTypes() []string {
	reg := stepRegistry()
	types := make([]string, 0, len(reg))
	for k := range reg {
		types = append(types, k)
	}
	return types
}
