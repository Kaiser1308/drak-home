#include "drak/id.hpp"

#include <string>

#include <gtest/gtest.h>

TEST(IdTest, AcceptsContractExampleIds) {
    // Values taken from contracts/v1/examples (command.json, manifest.json).
    EXPECT_TRUE(drak::is_valid_id("cmd-1"));
    EXPECT_TRUE(drak::is_valid_id("boot-001"));
    EXPECT_TRUE(drak::is_valid_id("room-01"));
    EXPECT_TRUE(drak::is_valid_id("a"));
    EXPECT_TRUE(drak::is_valid_id("A1_b.c-d"));
}

TEST(IdTest, RejectsBadFirstCharacter) {
    EXPECT_FALSE(drak::is_valid_id(""));
    EXPECT_FALSE(drak::is_valid_id("-cmd"));
    EXPECT_FALSE(drak::is_valid_id("_cmd"));
    EXPECT_FALSE(drak::is_valid_id(".cmd"));
}

TEST(IdTest, RejectsCharactersOutsidePattern) {
    EXPECT_FALSE(drak::is_valid_id("cmd 1"));
    EXPECT_FALSE(drak::is_valid_id("cmd/1"));
    EXPECT_FALSE(drak::is_valid_id("phòng-01"));
}

TEST(IdTest, EnforcesMaxLength) {
    EXPECT_TRUE(drak::is_valid_id(std::string(drak::kMaxIdLength, 'a')));
    EXPECT_FALSE(drak::is_valid_id(std::string(drak::kMaxIdLength + 1, 'a')));
}
