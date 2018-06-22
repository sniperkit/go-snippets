/* Copyright (C) 2016 {{cookiecutter.author_name}}
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

func main() {
	agent, err := New{{cookiecutter.agent_var}}()

	if err != nil {
		panic(err)
	}

	if err = agent.ParseCommandLine(); err != nil {
		panic(err)	
	}

	if err = agent.Wait(); err != nil {
		panic(err)	
	}

	if err = agent.Close(); err != nil {
		panic(err)	
	}
}
