#include <string>
#include <iostream>
#include <fstream>
#include <vector>
#include <algorithm>
#include <chrono>
#include <random>

// Edge List Storage
uint32_t num_vertices = 0;
std::vector< std::pair< uint32_t, uint32_t > > edge_list;

// CSR Storage
std::vector< uint64_t >  vert_offset;
std::vector< uint32_t >  edge_targets;


int main()
{
  //
  // Read the edges & build sorted edgelist
  {
    std::ifstream ifs("com-friendster.ungraph.txt");
    while(ifs && ifs.peek() == '#') { //skip comment lines
      std::string line;
      std::getline(ifs,line);
    }
    while(ifs) {
      uint32_t source,target;
      ifs >> source >> target;
      num_vertices = std::max(num_vertices, std::max(source-1,target-1));
      edge_list.push_back(std::make_pair(source-1,target-1));
      edge_list.push_back(std::make_pair(target-1,source-1)); //insert symmetric
    }
    num_vertices++;
    std::sort(edge_list.begin(), edge_list.end());
  }

  std::cout << "Read " << edge_list.size() << " edges, # vertices = " << num_vertices << std::endl;

  //
  //  Build CSR
  {
    std::vector<uint64_t> vertex_degree(num_vertices,0);
    edge_targets.reserve(edge_list.size());
    for(auto& e : edge_list) {
      vertex_degree[e.first]++;
      edge_targets.push_back(e.second);
    }
    vert_offset.reserve(num_vertices+1);
    vert_offset.push_back(0);
    uint64_t running_offset(0);
    for(auto degree : vertex_degree) {
      running_offset += degree;
      vert_offset.push_back(running_offset);
    }
  }

  //RNG for source selection
  std::mt19937 rng;    // random-number engine used (Mersenne-Twister in this case)
  std::uniform_int_distribution<uint32_t> uni(0,num_vertices-1); // guaranteed unbiased

  // Run 10 bfs trials
  double total_time(0);
  for(size_t i=0; i<10; ++i) {
    uint32_t source(0);
    do { //loops until a source with edges is found.
      source = 0;
    } while(vert_offset[source+1] - vert_offset[source] == 0);
    std::cout << "Source = " << source << std::endl;
    { // Start BFS
      auto wcts = std::chrono::system_clock::now();
      std::vector<uint8_t> vert_level(num_vertices, 255);
      std::vector<bool> vert_visited(num_vertices, false);
      uint8_t n_level = 2;
      std::vector<uint32_t> cur_level, next_level;
      cur_level.reserve(num_vertices); next_level.reserve(num_vertices);
      vert_level[source] = 0;
      vert_visited[source] = 1;
      cur_level.push_back(source);
      while(!cur_level.empty()) {
        for(auto v : cur_level) {
          for(size_t i=vert_offset[v]; i < vert_offset[v+1]; ++i) {
            uint32_t neighbor = edge_targets[i];
            if(!vert_visited[neighbor]) {
              next_level.push_back(neighbor);
              vert_level[neighbor] = n_level;
              vert_visited[neighbor] = true;
            }
          }
        }
        std::cout << "Completed level " << (int) n_level-1 << " size = " << next_level.size() << std::endl;
        ++n_level;
        cur_level.clear();
        next_level.swap(cur_level);
        std::sort(cur_level.begin(), cur_level.end());
      }
      std::chrono::duration<double> wctduration = (std::chrono::system_clock::now() - wcts);
      std::cout << "Finished in " << wctduration.count() << " seconds [Wall Clock]" << std::endl;
      total_time += wctduration.count();
    }
  }
  std::cout << "Average time for 10 runs = " << total_time / double(10) << std::endl;
  return 0;
}
