#version 410

uniform sampler2D tex;
uniform mat4 camera;
uniform mat4 model;
uniform vec3 camPosition;

in vec2 fragTexCoord;
in vec3 fragNormal;
in vec3 fragPosition;

out vec4 outputColor;

const float shininess = 0.5;
const vec3 lightColor = vec3(1, 1, 1);
const vec3 lightPosition = vec3(0,100,0);
const float lightPower = 40.0;
const vec3 ambientColor = vec3(0.4,0.4,0.4);
const vec3 diffuseColor = vec3(0.4, 0.4, 0.4);
const vec3 specColor = vec3(0.4,0.4,0.4);

void main() {
  vec3 color = texture(tex, vec2(fragTexCoord.x, 1.0-fragTexCoord.y)).rgb;
  vec3 normal = normalize(fragNormal);
  vec3 lightDir = normalize(lightPosition - fragPosition);
  
  float lambertian = max(dot(lightDir, normal), 0.0);
  float spec = 0.0;

  if (lambertian > 0) {
    vec3 viewDir = normalize(camPosition-fragPosition);
    vec3 halfDir = normalize(lightDir + viewDir);
    float specAngle = max(dot(halfDir, normal), 0);
    spec = pow(specAngle, shininess);
  }

  vec3 ambient = ambientColor * color;
  vec3 diffuse = lambertian * color;
  vec3 specular = specColor * spec;
  
  outputColor = vec4(max(diffuse, ambient), 1);
}
