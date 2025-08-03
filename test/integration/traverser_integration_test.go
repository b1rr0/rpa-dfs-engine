package integration

import (
	"os"
	"path/filepath"
	"testing"

	"rpa-dfs-engine/internal/traverser"
	"rpa-dfs-engine/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraverser_EndToEndWorkflows(t *testing.T) {
	// Setup temporary directory for test files
	tmpDir := t.TempDir()

	t.Run("Complete login workflow execution", func(t *testing.T) {
		// Create workflow file
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "{{user.website}}",
				"next": {
					"nodeType": "wait",
					"duration": 1000,
					"next": {
						"nodeType": "fillField",
						"selector": "{{user.emailSelector}}",
						"value": "{{user.email}}",
						"next": {
							"nodeType": "fillField",
							"selector": "{{user.passwordSelector}}",
							"value": "{{user.password}}",
							"next": {
								"nodeType": "clickButton",
								"selector": "{{user.loginButtonSelector}}",
								"next": null
							}
						}
					}
				}
			},
			"metadata": {
				"name": "Login Workflow",
				"version": "1.0.0",
				"description": "Complete login workflow with form filling"
			}
		}`

		// Create context file
		contextJSON := `{
			"user": {
				"website": "https://example.com/login",
				"email": "test@example.com",
				"password": "secretpassword123",
				"emailSelector": "#email",
				"passwordSelector": "#password",
				"loginButtonSelector": "#login-btn"
			}
		}`

		workflowFile := createTestFile(t, tmpDir, "login_workflow.json", workflowJSON)
		contextFile := createTestFile(t, tmpDir, "login_context.json", contextJSON)

		// Execute workflow with mock browser
		mockBrowser := mocks.NewMockTraverserBrowser()
		traverser := traverser.NewWithBrowser(mockBrowser)
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, contextFile)
		require.NoError(t, err)

		// Verify complete execution
		assert.True(t, mockBrowser.HasNavigatedTo("https://example.com/login"))
		assert.Equal(t, "test@example.com", mockBrowser.GetFieldValue("#email"))
		assert.Equal(t, "secretpassword123", mockBrowser.GetFieldValue("#password"))
		assert.True(t, mockBrowser.WasElementClicked("#login-btn"))
	})

	t.Run("File upload workflow with forEach", func(t *testing.T) {
		// Note: This test mocks the forEach interaction since we can't simulate user input in tests
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "{{user.uploadPage}}",
				"next": {
					"nodeType": "sequence",
					"sequence": [
						{
							"nodeType": "sendFile",
							"selector": "input[type='file']",
							"filePath": "{{user.files[0]}}"
						},
						{
							"nodeType": "wait",
							"duration": 1000
						},
						{
							"nodeType": "clickButton",
							"selector": "#upload-btn"
						}
					],
					"next": null
				}
			},
			"metadata": {
				"name": "File Upload Workflow",
				"version": "1.0.0"
			}
		}`

		contextJSON := `{
			"user": {
				"uploadPage": "https://example.com/upload",
				"files": ["/path/to/document1.pdf", "/path/to/image1.jpg"]
			}
		}`

		workflowFile := createTestFile(t, tmpDir, "upload_workflow.json", workflowJSON)
		contextFile := createTestFile(t, tmpDir, "upload_context.json", contextJSON)

		mockBrowser := mocks.NewMockTraverserBrowser()
		traverser := traverser.NewWithBrowser(mockBrowser)
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, contextFile)
		require.NoError(t, err)

		assert.True(t, mockBrowser.HasNavigatedTo("https://example.com/upload"))
		assert.Equal(t, "/path/to/document1.pdf", mockBrowser.GetUploadedFile("input[type='file']"))
		assert.True(t, mockBrowser.WasElementClicked("#upload-btn"))
	})

	t.Run("Conditional workflow execution", func(t *testing.T) {
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "{{user.website}}",
				"next": {
					"nodeType": "conditional",
					"conditionExpression": "{{user.age}} >= 18",
					"branches": {
						"yes": {
							"nodeType": "sequence",
							"sequence": [
								{
									"nodeType": "fillField",
									"selector": "#age-category",
									"value": "adult"
								},
								{
									"nodeType": "fillField",
									"selector": "#permissions",
									"value": "full"
								}
							],
							"next": null
						},
						"no": {
							"nodeType": "fillField",
							"selector": "#age-category",
							"value": "minor",
							"next": null
						}
					}
				}
			},
			"metadata": {
				"name": "Age-based Conditional Workflow",
				"version": "1.0.0"
			}
		}`

		// Test adult path
		contextJSONAdult := `{
			"user": {
				"website": "https://example.com/register",
				"age": 25
			}
		}`

		workflowFile := createTestFile(t, tmpDir, "conditional_workflow.json", workflowJSON)
		contextFileAdult := createTestFile(t, tmpDir, "adult_context.json", contextJSONAdult)

		mockBrowser := mocks.NewMockTraverserBrowser()
		traverser := traverser.NewWithBrowser(mockBrowser)
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, contextFileAdult)
		require.NoError(t, err)

		assert.Equal(t, "adult", mockBrowser.GetFieldValue("#age-category"))
		assert.Equal(t, "full", mockBrowser.GetFieldValue("#permissions"))

		// Test minor path
		contextJSONMinor := `{
			"user": {
				"website": "https://example.com/register",
				"age": 16
			}
		}`

		contextFileMinor := createTestFile(t, tmpDir, "minor_context.json", contextJSONMinor)
		mockBrowser.Reset()

		err = traverser.ExecuteWorkflow(workflowFile, contextFileMinor)
		require.NoError(t, err)

		assert.Equal(t, "minor", mockBrowser.GetFieldValue("#age-category"))
		assert.Equal(t, "", mockBrowser.GetFieldValue("#permissions")) // Should not be filled
	})

	t.Run("Complex nested workflow", func(t *testing.T) {
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "{{user.website}}",
				"next": {
					"nodeType": "fillField",
					"selector": "#username",
					"value": "{{user.username}}",
					"next": {
						"nodeType": "conditional",
						"conditionExpression": "{{user.hasPassword}} == true",
						"branches": {
							"yes": {
								"nodeType": "sequence",
								"sequence": [
									{
										"nodeType": "fillField",
										"selector": "#password",
										"value": "{{user.password}}"
									},
									{
										"nodeType": "fillField",
										"selector": "#confirm-password",
										"value": "{{user.password}}"
									}
								],
								"next": {
									"nodeType": "clickButton",
									"selector": "#submit",
									"next": null
								}
							},
							"no": {
								"nodeType": "clickButton",
								"selector": "#guest-login",
								"next": null
							}
						}
					}
				}
			},
			"metadata": {
				"name": "Complex Nested Workflow",
				"version": "1.0.0"
			}
		}`

		contextJSON := `{
			"user": {
				"website": "https://example.com/login",
				"username": "testuser",
				"password": "secret123",
				"hasPassword": true
			}
		}`

		workflowFile := createTestFile(t, tmpDir, "complex_workflow.json", workflowJSON)
		contextFile := createTestFile(t, tmpDir, "complex_context.json", contextJSON)

		mockBrowser := mocks.NewMockTraverserBrowser()
		traverser := traverser.NewWithBrowser(mockBrowser)
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, contextFile)
		require.NoError(t, err)

		// Verify execution path
		assert.True(t, mockBrowser.HasNavigatedTo("https://example.com/login"))
		assert.Equal(t, "testuser", mockBrowser.GetFieldValue("#username"))
		assert.Equal(t, "secret123", mockBrowser.GetFieldValue("#password"))
		assert.Equal(t, "secret123", mockBrowser.GetFieldValue("#confirm-password"))
		assert.True(t, mockBrowser.WasElementClicked("#submit"))
		assert.False(t, mockBrowser.WasElementClicked("#guest-login"))
	})
}

func TestTraverser_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Browser failure during workflow", func(t *testing.T) {
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "https://example.com",
				"next": {
					"nodeType": "fillField",
					"selector": "#email",
					"value": "test@example.com",
					"next": null
				}
			},
			"metadata": {
				"name": "Test Workflow",
				"version": "1.0.0"
			}
		}`

		contextJSON := `{"user": {}}`

		workflowFile := createTestFile(t, tmpDir, "workflow.json", workflowJSON)
		contextFile := createTestFile(t, tmpDir, "context.json", contextJSON)

		// Configure mock browser to fail on fillField
		mockBrowser := mocks.NewMockTraverserBrowser()
		mockBrowser.SetFailureFor("fillField", "element not found")

		traverser := traverser.NewWithBrowser(mockBrowser)
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, contextFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fillField failed")

		// Navigation should still have succeeded
		assert.True(t, mockBrowser.HasNavigatedTo("https://example.com"))
	})

	t.Run("Missing context file", func(t *testing.T) {
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "https://example.com"
			},
			"metadata": {
				"name": "Test Workflow",
				"version": "1.0.0"
			}
		}`

		workflowFile := createTestFile(t, tmpDir, "workflow.json", workflowJSON)
		missingContextFile := filepath.Join(tmpDir, "missing_context.json")

		traverser := traverser.New()
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, missingContextFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load context")
	})
}

func TestTraverser_TemplateResolution(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Complex template resolution", func(t *testing.T) {
		workflowJSON := `{
			"graph": {
				"nodeType": "moveToPage",
				"url": "{{user.baseUrl}}/{{user.section}}",
				"next": {
					"nodeType": "fillField",
					"selector": "{{user.selectors.email}}",
					"value": "{{user.profile.firstName}} {{user.profile.lastName}} <{{user.profile.email}}>",
					"next": null
				}
			},
			"metadata": {
				"name": "Template Test Workflow",
				"version": "1.0.0"
			}
		}`

		contextJSON := `{
			"user": {
				"baseUrl": "https://example.com",
				"section": "profile",
				"selectors": {
					"email": "#user-email"
				},
				"profile": {
					"firstName": "John",
					"lastName": "Doe",
					"email": "john.doe@example.com"
				}
			}
		}`

		workflowFile := createTestFile(t, tmpDir, "template_workflow.json", workflowJSON)
		contextFile := createTestFile(t, tmpDir, "template_context.json", contextJSON)

		mockBrowser := mocks.NewMockTraverserBrowser()
		traverser := traverser.NewWithBrowser(mockBrowser)
		defer traverser.Close()

		err := traverser.ExecuteWorkflow(workflowFile, contextFile)
		require.NoError(t, err)

		// Verify template resolution
		assert.True(t, mockBrowser.HasNavigatedTo("https://example.com/profile"))
		assert.Equal(t, "John Doe <john.doe@example.com>", mockBrowser.GetFieldValue("#user-email"))
	})
}

// Helper function to create test files
func createTestFile(t *testing.T, dir, name, content string) string {
	filePath := filepath.Join(dir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
	return filePath
}
