#version 410

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;
in vec4 position;

void main() {
	if(gl_FrontFacing)
    	outputColor = texture(tex, fragTexCoord);
	else
		outputColor = vec4(1.0, 0.5,0.0,0.0);
	//outputColor = vec4(position.x*0.25, position.y*0.75, position.z, 0);
}
