class Chaosmonkey < Formula
  desc     "Controlled failure injection with Chaos Monkey"
  homepage "https://github.com/mlafeldt/chaosmonkey"
  version  "%VERSION%"
  url      "https://github.com/mlafeldt/chaosmonkey/releases/download/v#{version}/chaosmonkey_darwin_amd64"
  sha256   "%SHA%"

  bottle :unneeded

  def install
    bin.install "chaosmonkey_darwin_amd64" => "chaosmonkey"
  end

  test do
    system "#{bin}/chaosmonkey --version"
  end
end
