#version 410

// Input attributes
layout(location = 0) in vec2 fragTexCoord;
layout(location = 1) in vec4 position;

// Uniforms
uniform sampler2D tex;

// Output attributes
layout(location = 0) out vec4 outputColor;

void main() {
	if(gl_FrontFacing)
    	outputColor = texture(tex, fragTexCoord);
	else
		outputColor = vec4(1.0, 0.5,0.0,0.0);
	//outputColor = vec4(position.x*0.25, position.y*0.75, position.z, 0);
}
