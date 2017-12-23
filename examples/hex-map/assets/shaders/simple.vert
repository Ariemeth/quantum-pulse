#version 410

// Input attributes
layout(location = 0) in vec3 vert;
layout(location = 1) in vec2 vertTexCoord;

// Uniforms
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

// Output attributes
layout(location = 0) out vec2 fragTexCoord;
layout(location = 1) out vec4 position;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
	position = gl_Position;
}
