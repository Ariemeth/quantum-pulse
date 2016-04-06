#version 440

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;
in vec4 position;

void main() {
    //outputColor = texture(tex, fragTexCoord);
	outputColor = vec4(position.x*0.25, position.y*0.75, position.z, 0);
}
