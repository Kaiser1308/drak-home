# Pure C++ core shared by firmware and host tests. No ESP-IDF headers here.
idf_component_register(SRCS "src/id.cpp"
                       INCLUDE_DIRS "include")
