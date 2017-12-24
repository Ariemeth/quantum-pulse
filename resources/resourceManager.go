package resources

// Manager manages engine resources.
type Manager struct {
	sm ShaderManager
	tm TextureManager
}

// NewManager creates a new Manager to handle shader and texture assets.
func NewManager() *Manager {
	am := Manager{
		sm: newShaderManager(),
		tm: newTextureManager(),
	}
	return &am
}

// Shaders retrieves the ShaderManager.
func (am *Manager) Shaders() ShaderManager {
	return am.sm
}

// Textures retrieves the TextureManager.
func (am *Manager) Textures() TextureManager {
	return am.tm
}
