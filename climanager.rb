# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Climanager < Formula
  desc ""
  homepage "https://github.com/darrenleak/climanager"
  version "0.1.7"

  on_macos do
    url "https://github.com/darrenleak/climanager/releases/download/v0.1.7/climanager_0.1.7_darwin_all.tar.gz"
    sha256 "489ccf479cc27bbe00076a17e1745180c9b0d70b0f2f31ee55f968b8ac9a2c60"

    def install
      bin.install "climanager"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/darrenleak/climanager/releases/download/v0.1.7/climanager_0.1.7_linux_arm64.tar.gz"
      sha256 "c169171cbf92c91052d0c23fadf5110c623283c46beb3e3bcb66db8eef0922ad"

      def install
        bin.install "climanager"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/darrenleak/climanager/releases/download/v0.1.7/climanager_0.1.7_linux_amd64.tar.gz"
      sha256 "35999cbf5a89231fc6c75aac5c6949b35525cd73c602cc06ff7838c3a51343b4"

      def install
        bin.install "climanager"
      end
    end
  end
end
